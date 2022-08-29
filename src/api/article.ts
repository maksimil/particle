import type { MarkdownInstance } from "astro";
import * as path from "path";

export type IndexFrontmatter = {
  title: string;
  author: string;
  chapters: string[];
};

export type ChapterFrontmatter = {
  title: string;
};

export type Chapter = {
  id: string;
  title: string;
  md: MarkdownInstance<ChapterFrontmatter>;
};

export type Article = {
  id: string;
  title: string;
  author: string;
  chapters: Chapter[];
  md: MarkdownInstance<IndexFrontmatter>;
};

export const getArticles = async (): Promise<Article[]> => {
  //@ts-ignore
  const indexes = (await import.meta.glob(
    "/data/articles/*/index.md"
  )) as Record<string, () => Promise<MarkdownInstance<IndexFrontmatter>>>;

  const articles: Article[] = await Promise.all(
    Object.values(indexes).map(async (index_promise) => {
      const index = await index_promise();
      const id = path.basename(path.dirname(index.file));
      const chapters: Chapter[] = await Promise.all(
        index.frontmatter.chapters.map(async (chname) => {
          const chmd: MarkdownInstance<ChapterFrontmatter> = await import(
            `../../data/articles/${id}/${chname}.md`
          );
          return {
            id: chname,
            title: chmd.frontmatter.title,
            md: chmd,
          };
        })
      );
      return {
        id,
        title: index.frontmatter.title,
        author: index.frontmatter.author,
        chapters,
        md: index,
      };
    })
  );

  return articles;
};
