import * as path from "node:path";
import {
  Component,
  createEffect,
  createSignal,
  For,
  createMemo,
  Show,
} from "solid-js";
import { createStore, unwrap } from "solid-js/store";

const TextField: Component<{
  label: string;
  text: string;
  setText: (text: string) => void;
}> = (props) => {
  return (
    <div class="w-full flex mt-2">
      <span class="text-xl text-c3 mr-4 font-mono">{props.label}</span>
      <input
        class={
          "text-xl flex-1 bg-transparent focus:outline-none " +
          "font-mono text-c3 placeholder-c2 " +
          "border-b-1 border-c2 "
        }
        type="text"
        placeholder="..."
        value={props.text}
        onchange={(e) => {
          props.setText(e.currentTarget.value);
        }}
      />
    </div>
  );
};

const FileField: Component<{
  label: string;
  load: (text: string, name: string) => void;
}> = (props) => {
  const [fileName, setFileName] = createSignal("Upload file");
  const upload = () => {
    const fileSelector = document.createElement("input");
    fileSelector.setAttribute("type", "file");
    fileSelector.click();
    fileSelector.onchange = (_) => {
      if (fileSelector.files === null) {
        return;
      }
      const file = fileSelector.files[0];
      file.text().then((value) => {
        props.load(value, file.name);
        setFileName(file.name);
      });
    };
  };
  return (
    <div class="w-100 flex mt-2">
      <span class="text-xl text-c3 mr-4 font-mono">{props.label}</span>
      <button
        class="text-xl flex-1 font-mono text-c3 hover:text-c2 text-left"
        onclick={upload}
      >
        {fileName()}
      </button>
    </div>
  );
};

type SubmitData = {
  title: string;
  author: string;
  description: string | null;
  chapters: { title: string; contents: string | null; id: string }[];
};

const SubmitPage: Component<{}> = () => {
  const [store, setStore] = createStore<SubmitData>({
    title: "",
    author: "",
    description: null,
    chapters: [],
  });

  const [submitStatus, setSubmitStatus] = createSignal<
    | { ty: "not-sent" }
    | { ty: "sending" }
    | { ty: "success" }
    | { ty: "err"; code: number }
  >({ ty: "not-sent" }, { equals: false });
  const statusBar = createMemo(() => {
    switch (submitStatus().ty) {
      case "not-sent":
        return "";
      case "sending":
        return "Sending >>>";
      case "success":
        return "Successful <3";
      case "err":
        //@ts-ignore
        return `Failed with code ${submitStatus().code} :(`;
    }
  });

  createEffect(() => console.log(store));
  return (
    <div class="mt-2 w-100">
      <TextField
        label="title"
        text={store.title}
        setText={(text) => setStore("title", text)}
      />
      <TextField
        label="author"
        text={store.author}
        setText={(text) => setStore("author", text)}
      />
      <FileField
        label="description"
        load={(text) => setStore("description", text)}
      />
      <div class="text-xl text-c3 font-mono mt-2">chapters</div>
      <div class="w-100 space-y-5">
        <For each={store.chapters}>
          {(chapter, idx) => (
            <div class="border-l-1 border-c2 pl-2 ml-3">
              <TextField
                label="title"
                text={chapter.title}
                setText={(text) => setStore("chapters", idx(), "title", text)}
              />
              <FileField
                label="contents"
                load={(text, name) => {
                  setStore("chapters", idx(), "contents", text);
                  setStore("chapters", idx(), "id", name.split(".")[0]);
                }}
              />
              <button
                class="text-xl text-c3 hover:text-c2 font-mono mt-2"
                onclick={() => {
                  setStore("chapters", (c) => {
                    const cp = c.slice(0, idx());
                    const cn = c.slice(idx() + 1);
                    return [...cp, ...cn];
                  });
                }}
              >
                remove
              </button>
            </div>
          )}
        </For>
      </div>
      <button
        class="text-xl text-c3 hover:text-c2 font-mono mt-4"
        onclick={() => {
          setStore("chapters", (c) => [
            ...c,
            { title: "", contents: null, id: "none" },
          ]);
        }}
      >
        add chapter
      </button>
      <br />
      <button
        class="text-xl text-c3 hover:text-c2 font-mono mt-4"
        onclick={() => {
          setSubmitStatus({ ty: "sending" });
          const data = unwrap(store);
          console.log("Submitting", data);
          const body = JSON.stringify(data);
          fetch("/api/submit", {
            method: "POST",
            mode: "cors",
            cache: "no-cache",
            credentials: "same-origin",
            headers: {
              "Content-Type": "application/json",
            },
            redirect: "follow",
            referrerPolicy: "no-referrer",
            body,
          }).then((v) => {
            if (v.status === 200) {
              setSubmitStatus({ ty: "success" });
            } else {
              setSubmitStatus({ ty: "err", code: v.status });
            }
          });
        }}
      >
        submit ^_^
      </button>
      <Show when={statusBar() !== ""}>
        <div class="text-xl text-c3 font-mono mt-4">{statusBar()}</div>
      </Show>
    </div>
  );
};

export default SubmitPage;
