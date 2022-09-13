import {
  Component,
  createSignal,
  Switch,
  Match,
  onMount,
  onCleanup,
} from "solid-js";

const LinkButton: Component<{ label: string; url: string }> = (props) => {
  return (
    <a
      href={props.url}
      class={
        "text-3xl flex-1 font-mono tracking-wide text-center " +
        "text-white hover:text-white/50 z-5 "
      }
    >
      {props.label}
    </a>
  );
};

const Particle: Component = () => {
  const [enter, setEnter] = createSignal(false);

  let hitDiv: HTMLDivElement;
  let screenDiv: HTMLDivElement;

  onMount(() => {
    const mouseCheck = (e: MouseEvent) => {
      console.log("check");
      setEnter(false);
      const rect = hitDiv.getBoundingClientRect();
      setEnter(
        rect.left <= e.clientX &&
          e.clientX <= rect.right &&
          rect.top <= e.clientY &&
          e.clientY <= rect.bottom
      );
    };

    screenDiv.addEventListener("mousemove", mouseCheck);
    onCleanup(() => {
      screenDiv.removeEventListener("mousemove", mouseCheck);
    });
  });

  return (
    <div
      ref={screenDiv}
      class="w-screen h-screen flex justify-center items-center bg-c3"
    >
      <div class="w-300px h-116px">
        <div ref={hitDiv} class="w-300px h-116px absolute " />
        <Switch>
          <Match when={enter()}>
            <div class="h-88px flex flex-row justify-center flex-wrap">
              <LinkButton label="browse" url="/browse" />
              <LinkButton label="submit" url="/submit" />
              <LinkButton
                label="github<3"
                url="https://github.com/maksimil/particle"
              />
            </div>
          </Match>
          <Match when={!enter()}>
            <div class="h-88px text-5xl tracking-widest text-center text-white">
              Particle.
            </div>
          </Match>
        </Switch>
        <div class="text-xl text-center text-white">
          {"A place for stuff <3"}
        </div>
      </div>
    </div>
  );
};

export default Particle;
