import { defineCollection, z } from 'astro:content';

const blog = defineCollection({
  type: 'content',
  schema: z.object({
    title:       z.string(),
    date:        z.coerce.date(),
    tags:        z.array(z.string()).default([]),
    excerpt:     z.string().optional(),
    draft:       z.boolean().default(false),
    readTime:    z.string().optional(),
  }),
});

export const collections = { blog };
