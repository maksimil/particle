import { Component, createSignal } from "solid-js";

const SearchBrowse: Component<{ links: [string, string][] }> = (props) => {
  const [search, setSearch] = createSignal("");
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
        {props.links
          .filter(([name, _]) => name.startsWith(search()))
          .map(([name, link]) => {
            return (
              <div>
                <a href={link} class="text-xl font-mono text-c3 hover:text-c2">
                  <span class="text-c2">
                    {name.substring(0, search().length)}
                  </span>
                  {name.substring(search().length, name.length)}
                </a>
              </div>
            );
          })}
      </div>
    </>
  );
};

export default SearchBrowse;
