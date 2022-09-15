import { Component, createSignal } from "solid-js";
import type { Article } from "@api/article";

type SearchItem = {
  searchValue: string;
  link: string;
  isTitle: boolean;
  Component: Component<{ search: string }>;
};

const toSearchItems = (article: Article): SearchItem[] => {
  const titleItem: SearchItem = {
    searchValue: article.title.toLowerCase(),
    link: `/article/${article.id}`,
    isTitle: true,
    Component: (props) => {
      return (
        <>
          <span class="text-c2">
            {article.title.substring(0, props.search.length)}
          </span>
          {article.title.substring(props.search.length, article.title.length)}
        </>
      );
    },
  };

  const chapterItems: SearchItem[] = article.chapters.map(
    (chapter): SearchItem => {
      return {
        searchValue: chapter.title.toLowerCase(),
        link: `/article/${article.id}/${chapter.id}`,
        isTitle: false,
        Component: (props) => {
          return (
            <>
              <span class="text-c2">
                {`${article.title} -> ${chapter.title.substring(
                  0,
                  props.search.length
                )}`}
              </span>
              {chapter.title.substring(
                props.search.length,
                chapter.title.length
              )}
            </>
          );
        },
      };
    }
  );

  return [titleItem, ...chapterItems];
};

const SearchBrowse: Component<{ articles: Article[] }> = (props) => {
  const [search, setSearch] = createSignal("");

  const items = props.articles.map(toSearchItems).flat();

  return (
    <>
      <div>
        <input
          class={
            "text-xl bg-transparent focus:outline-none " +
            "font-mono text-c2 placeholder-c2 "
          }
          type="text"
          placeholder="Search"
          value={search()}
          oninput={(e) => {
            setSearch(e.currentTarget.value);
          }}
        />
      </div>
      <div class="overflow-x-scroll flex-1">
        {items
          .filter((item) => {
            if (search() === "") {
              return item.isTitle;
            } else {
              return item.searchValue.startsWith(search().toLowerCase());
            }
          })
          .map((item) => {
            return (
              <div>
                <a
                  href={item.link}
                  class="text-xl font-mono text-c3 hover:text-c2"
                >
                  {item.Component({ search: search() })}
                </a>
              </div>
            );
          })}
      </div>
    </>
  );
};

export default SearchBrowse;
