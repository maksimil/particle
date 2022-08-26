import * as yaml from "js-yaml";
import * as path from "path";
import type { MarkdownInstance } from "astro";

export type ArticleMeta = {
  title: string;
  description?: string;
  author?: string;
  chapters: string[];
};

export type Article = {
  title: string;
  id: string;
  description?: string;
  author?: string;
  chapters: Chapter[];
};

export type Chapter = {
  title: string;
  contents: string;
};

const getArticles = async (): Promise<Article[]> => {
  //@ts-ignore
  const indexfiles = await import.meta.glob("../../data/articles/*/index.yml");

  return await Promise.all(
    Object.keys(indexfiles).map(async (fname) => {
      const { default: yamlraw } = await import(
        /* @vite-ignore */
        fname + "?raw"
      );
      const metadata = yaml.load(yamlraw) as ArticleMeta;
      const chapters = await Promise.all(
        metadata.chapters.map(async (mdname) => {
          const mddata: MarkdownInstance<{ title: string }> = await import(
            /* @vite-ignore */
            path.join(path.dirname(fname), mdname)
          );

          return {
            title: mddata.frontmatter.title,
            contents: mddata.compiledContent(),
          } as Chapter;
        })
      );

      return {
        title: metadata.title,
        id: path.basename(path.dirname(fname)),
        description: metadata.description,
        author: metadata.author,
        chapters,
      };
    })
  );
};

export const ARTICLES = await getArticles();

export const getArticle = async (id: string): Promise<Article> => {
  const fd = ARTICLES.find((v) => v.id === id);

  if (fd === undefined) {
    throw Error(`Article with id:${id} was not found`);
  }

  return fd;
};
