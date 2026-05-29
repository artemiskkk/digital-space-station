---
title: "ETA5003 零穿越检测：从硬件信号到软件双阈值决策"
date: 2025-07-18
tags: ["嵌入式", "PLC", "GD32", "C++"]
excerpt: "在低压电力线载波通信场景下，工频干扰导致过零点检测误触发是个棘手问题。本文记录了基于 ETA5003 的软件双阈值方案，以及在三相系统中交叉验证的工程实践。"
readTime: "8 min read"
draft: false
---

在低压电力线载波 (PLC) 系统中，CCO 节点需要精确感知电网的**相位信息**，才能实现帧同步与区域识别。然而工频环境远不像教科书那样干净——谐波、开关电源噪声、感性负载的反射都会在过零点附近产生毛刺，导致频繁的误触发。

## 硬件拓扑回顾

ETA5003 是一款专为电力线通信设计的模拟前端 IC，其零穿越检测引脚 `ZCD` 在信号穿越阈值时输出一个脉冲。我们将它接入 GD32F103 的 `PA8`（EXTI Line 8），配置为**双边沿触发**。

```
AC 220V ──[变压器]──[ETA5003 ZCD]──[PC817 光耦]──[GD32 PA8/EXTI8]
                          │
                    [隔离电源 3.3V]
```

> **注意**：光耦会引入约 5～15μs 的传播延迟，必须在软件中补偿，否则相位计算会有系统性偏差。

## 核心问题：为什么单阈值不够？

单阈值检测假设信号在过零点附近是单调的。现实中，一个典型的"毛刺触发"场景：

1. 信号从 +ε 跌到 -ε（正常过零），触发一次中断
2. 叠加的噪声让信号在 -ε 附近来回震荡，再触发 2～3 次
3. 相位计算器看到了"每工频周期 4～5 次过零"，完全错乱

---

## 双阈值决策算法（C++ 实现）

核心思想：过零事件必须满足**两个连续采样点分别处于阈值带的两侧**，才判定为有效过零。

```cpp
// zero_crossing_detector.h
#pragma once
#include <cstdint>

namespace PLC {

/// 双阈值零穿越检测器
/// 采样率: 10 kHz (定时器每 100μs 调用一次 Update)
class ZeroCrossingDetector {
public:
    static constexpr int16_t kUpperThreshold =  512;  // ADC counts, ~10% Vref
    static constexpr int16_t kLowerThreshold = -512;
    static constexpr uint32_t kDebounceUs    = 800;   // 最小间隔 0.8ms (~4% 周期)

    struct Event {
        uint32_t timestamp_us;   ///< 捕获时刻 (微秒, 来自 DWT->CYCCNT)
        bool     rising;         ///< true = 负→正穿越
    };

    /// 喂入新的 ADC 采样值, 若检测到有效过零则返回 true
    bool Update(int16_t sample, uint32_t now_us, Event* out) {
        bool detected = false;

        if (prev_state_ == State::kAboveUpper && sample < kLowerThreshold) {
            // 正→负有效穿越（信号已明确越过下阈值）
            if (now_us - last_event_us_ > kDebounceUs) {
                out->timestamp_us = now_us;
                out->rising       = false;
                last_event_us_    = now_us;
                detected          = true;
            }
            prev_state_ = State::kBelowLower;
        } else if (prev_state_ == State::kBelowLower && sample > kUpperThreshold) {
            // 负→正有效穿越
            if (now_us - last_event_us_ > kDebounceUs) {
                out->timestamp_us = now_us;
                out->rising       = true;
                last_event_us_    = now_us;
                detected          = true;
            }
            prev_state_ = State::kAboveUpper;
        }
        // kInBand 状态：信号在阈值带内，不更新 prev_state_

        return detected;
    }

private:
    enum class State { kAboveUpper, kInBand, kBelowLower };
    State    prev_state_    = State::kInBand;
    uint32_t last_event_us_ = 0;
};

} // namespace PLC
```

## 三相交叉验证

单相检测仍有盲区。当三相电都接入系统时，可以利用**相位互差 120°** 这一约束进行交叉校验：

```go
// cross_validator.go  (上位机验证逻辑，用 Go 编写)
package plc

import "math"

const expectedDeltaDeg = 120.0

// ValidateThreePhase 检验三路时间戳是否符合 120° 间隔
// timestamps: [phaseA, phaseB, phaseC] 单位 μs，同一工频周期内的过零时刻
func ValidateThreePhase(timestamps [3]uint32, periodUs uint32) bool {
    period := float64(periodUs)

    for i := 0; i < 3; i++ {
        j := (i + 1) % 3
        deltaUs  := float64(timestamps[j]-timestamps[i])
        deltaDeg := deltaUs / period * 360.0

        // 允许 ±5° 的误差窗口
        if math.Abs(deltaDeg-expectedDeltaDeg) > 5.0 {
            return false
        }
    }
    return true
}
```

## 工程经验小结

- **光耦延迟补偿**：在软件中将所有时间戳减去一个固定偏置（实测约 `12μs`），或在标定阶段自动计算。
- **阈值自适应**：若电网电压波动大，可将 `kUpperThreshold` / `kLowerThreshold` 改为 ADC 满量程的百分比，而不是固定值。
- **GD32 EXTI 配置**：记得在 `EXTI_InitTypeDef` 中将 `EXTI_Trigger` 设为 `EXTI_Trigger_Rising_Falling`，同时使能 `AFIO` 时钟，否则引脚复用不生效——这个坑我踩过两次。

> 关于 ETA5003 的完整寄存器配置，可以参考官方应用笔记 AN-2301，里面的 ZCD 滤波器参数建议在实际硬件上重新标定，不要直接照搬默认值。

---

下一篇我打算写三相系统下的**区域识别算法**，以及如何利用过零相位差来区分同相/跨相站点。敬请期待。
