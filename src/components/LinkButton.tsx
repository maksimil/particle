import type { Component } from "solid-js";

const LinkButton: Component<{ url: string }> = (props) => {
  let imgRef: HTMLImageElement;

  return (
    <button class="w-9 h-9">
      <img
        ref={
          //@ts-ignore
          imgRef
        }
        src="/link.svg"
        class="w-full h-full"
        onmouseover={() => {
          imgRef.src = "/link-hover.svg";
        }}
        onmouseout={() => {
          imgRef.src = "/link.svg";
        }}
        onmousedown={() => {
          imgRef.src = "/link-down.svg";
        }}
        onmouseup={() => {
          imgRef.src = "/link-hover.svg";
        }}
        onclick={() => {
          console.log(props.url);
          navigator.clipboard.writeText(window.location.hostname + props.url);
        }}
      />
    </button>
  );
};

export default LinkButton;
