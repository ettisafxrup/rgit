// Shared script for index.html and docs.html.
// Each feature returns early when its elements are not on the page.

const RELEASE = "https://github.com/ettisafxrup/rgit/releases/latest/download/";
const reduceMotion = matchMedia("(prefers-reduced-motion: reduce)").matches;

const $ = (selector, root = document) => root.querySelector(selector);
const $$ = (selector, root = document) => [...root.querySelectorAll(selector)];
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

// Calls callback(element) the first time each element scrolls into view.
function whenVisible(elements, callback) {
  const observer = new IntersectionObserver((entries) => {
    for (const entry of entries) {
      if (entry.isIntersecting) {
        observer.unobserve(entry.target);
        callback(entry.target);
      }
    }
  }, { rootMargin: "0px 0px -10% 0px" });
  elements.forEach((el) => observer.observe(el));
}

function setupNav() {
  const nav = $(".nav");
  const progress = $(".progress");
  const toTop = $(".to-top");

  const onScroll = () => {
    const scrollable = document.documentElement.scrollHeight - innerHeight;
    nav.classList.toggle("scrolled", scrollY > 8);
    progress.style.transform = `scaleX(${scrollable > 0 ? scrollY / scrollable : 0})`;
    toTop.classList.toggle("show", scrollY > 900);
  };
  addEventListener("scroll", onScroll, { passive: true });
  onScroll();

  toTop.addEventListener("click", () => scrollTo({ top: 0 }));

  const menuButton = $(".menu-btn");
  const menu = $(".mobile-menu");
  const setMenu = (open) => {
    menu.hidden = !open;
    menuButton.setAttribute("aria-expanded", open);
  };
  menuButton.addEventListener("click", () => setMenu(menu.hidden));
  $$("a", menu).forEach((link) => link.addEventListener("click", () => setMenu(false)));
}

function setupReveal() {
  $$("[data-stagger]").forEach((group) => {
    $$("[data-reveal]", group).forEach((el, i) => el.style.setProperty("--delay", `${(i % 8) * 70}ms`));
  });
  whenVisible($$("[data-reveal]"), (el) => el.classList.add("in"));
}

function setupCounters() {
  whenVisible($$("[data-count]"), async (el) => {
    const target = Number(el.dataset.count);
    if (reduceMotion) return;
    for (let step = 1; step <= 30; step++) {
      el.textContent = Math.round(target * (1 - (1 - step / 30) ** 3));
      await sleep(35);
    }
  });
}

function setupCopyButtons() {
  for (const block of $$(".code")) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "copy";
    button.ariaLabel = "Copy";
    button.innerHTML = '<svg aria-hidden="true"><use href="#i-copy"/></svg>';
    button.addEventListener("click", async () => {
      try {
        await navigator.clipboard.writeText($("pre", block).textContent.trim());
      } catch {
        return; // clipboard access was blocked by the browser
      }
      button.classList.add("done");
      $("use", button).setAttribute("href", "#i-check");
      await sleep(1500);
      button.classList.remove("done");
      $("use", button).setAttribute("href", "#i-copy");
    });
    block.append(button);
  }
}

function selectTab(tab) {
  for (const other of $$('[role="tab"]', tab.parentElement)) {
    const selected = other === tab;
    other.setAttribute("aria-selected", selected);
    other.tabIndex = selected ? 0 : -1;
    document.getElementById(other.getAttribute("aria-controls")).hidden = !selected;
  }
}

function setupTabs() {
  for (const tab of $$('[role="tab"]')) {
    tab.addEventListener("click", () => selectTab(tab));
    tab.addEventListener("keydown", (event) => {
      const tabs = $$('[role="tab"]', tab.parentElement);
      const step = { ArrowRight: 1, ArrowLeft: -1 }[event.key];
      if (!step) return;
      const next = tabs[(tabs.indexOf(tab) + step + tabs.length) % tabs.length];
      selectTab(next);
      next.focus();
    });
  }
}

// ---------------------------------------------------------------- downloads

const DOWNLOADS = {
  windows: { label: "Download for Windows", file: "rgit-windows-x64-setup.exe", detail: "Installer · Windows 10 & 11" },
  "macos-arm64": { label: "Download for macOS", file: "rgit-macos-arm64", detail: "Apple Silicon (M1 and later)" },
  "macos-x64": { label: "Download for macOS", file: "rgit-macos-x64", detail: "Intel Mac" },
  "linux-x64": { label: "Download for Linux", file: "rgit-linux-x64", detail: "x64 · any distribution" },
  "linux-arm64": { label: "Download for Linux", file: "rgit-linux-arm64", detail: "ARM64 · any distribution" },
};

async function detectSystem() {
  const text = `${navigator.userAgentData?.platform ?? ""} ${navigator.userAgent}`.toLowerCase();
  if (/android|iphone|ipad|mobile/.test(text)) return null;

  let os = null;
  if (text.includes("win")) os = "windows";
  else if (text.includes("mac")) os = "macos";
  else if (text.includes("linux") || text.includes("x11")) os = "linux";
  if (!os) return null;

  // Most Macs today use Apple Silicon; Chromium browsers can confirm the CPU.
  let arch = os === "macos" || /aarch64|arm64/.test(text) ? "arm64" : "x64";
  try {
    const { architecture } = await navigator.userAgentData.getHighEntropyValues(["architecture"]);
    if (architecture === "arm") arch = "arm64";
    if (architecture === "x86") arch = "x64";
  } catch {
    // Not available in this browser: keep the guess.
  }
  return { os, key: os === "windows" ? "windows" : `${os}-${arch}` };
}

async function setupDownloads() {
  const system = await detectSystem();
  if (!system) return;

  const choice = DOWNLOADS[system.key];
  const button = $("#get-btn");
  if (button) {
    button.href = RELEASE + choice.file;
    $("#get-label").textContent = choice.label;
    $("#get-file").innerHTML = `<strong>${choice.detail}</strong>${choice.file}`;
    $("#get-cmd").hidden = system.os === "windows";
  }

  const tab = $(`.os-tab[data-os="${system.os}"]`);
  if (tab) selectTab(tab);
}

// ---------------------------------------------------------------- typing terminal

const PROMPT = '<span class="t-prompt">$</span> ';

const EXAMPLES = [
  {
    title: "~/projects/webapp — push",
    command: 'rgit push "Add dark mode"',
    output: [
      ' <span class="t-step">•</span> 3 changed files:',
      '   <span class="t-mod">modified  </span> src/theme.css',
      '   <span class="t-mod">modified  </span> src/app.js',
      '   <span class="t-new">new       </span> src/dark.css',
      ' <span class="t-ask">?</span> Stage all of these changes? <span class="t-dim">[Y/n]</span> y',
      ' <span class="t-ok">✓</span> Committed <span class="t-hash">4f1e02f</span> Add dark mode',
      ' <span class="t-step">›</span> Pulling latest changes from <b>origin/main</b>',
      ' <span class="t-step">›</span> Pushing <b>main</b>',
      ' <span class="t-ok">✓</span> Pushed <b>main</b> to <b>origin/main</b>',
    ],
  },
  {
    title: "~/projects/webapp — status",
    command: "rgit status",
    output: [
      "",
      '   <span class="t-dim">Branch     </span>  <b>search</b>  → origin/search  <span class="t-mod">2 to push</span>',
      '   <span class="t-dim">Last commit</span>  <span class="t-hash">9b80fc3</span> Add search box <span class="t-dim">(3 minutes ago)</span>',
      '   <span class="t-dim">Remote     </span>  origin <span class="t-dim">git@github.com:you/webapp.git</span>',
      '   <span class="t-dim">Changes    </span>  <span class="t-mod">1 modified</span>, <span class="t-step">1 new</span>',
      '   <span class="t-mod">modified  </span> src/search.js',
      '   <span class="t-new">new       </span> src/search.css',
    ],
  },
  {
    title: "~/projects/webapp — undo",
    command: "rgit undo",
    output: [
      ' <span class="t-step">•</span> Last commit: <span class="t-hash">9b80fc3</span> Add search box <span class="t-dim">(2 minutes ago)</span>',
      ' <span class="t-ask">?</span> Undo this commit? Its changes will stay staged. <span class="t-dim">[Y/n]</span> y',
      ' <span class="t-ok">✓</span> Undid <span class="t-hash">9b80fc3</span> — your changes are still staged.',
    ],
  },
  {
    title: "~/projects/api — new project files",
    command: "rgit ignore node vscode && rgit license mit",
    output: [
      ' <span class="t-ok">✓</span> Added 22 entries to <b>.gitignore</b>',
      ' <span class="t-ask">?</span> Copyright year: <span class="t-dim">(2026)</span>',
      ' <span class="t-ask">?</span> Copyright holder (your name or organization): <span class="t-dim">(Jane Doe)</span>',
      ' <span class="t-ok">✓</span> Created <b>LICENSE</b> with the MIT License',
    ],
  },
  {
    title: "~/projects/api — typo",
    command: "rgit psuh",
    output: [' <span class="t-err">✗</span> unknown command "psuh" — did you mean "push"? (see \'rgit help\')'],
  },
];

function setupTypingTerminal() {
  const terminal = $("#typer");
  if (!terminal) return;
  const screen = $("pre", terminal);
  const title = $(".title", terminal);
  const dots = EXAMPLES.map((example, i) => {
    const dot = document.createElement("button");
    dot.type = "button";
    dot.ariaLabel = `Show example: ${example.command}`;
    dot.addEventListener("click", () => play(i));
    $(".typer-tabs", terminal).append(dot);
    return dot;
  });

  let onScreen = true;
  new IntersectionObserver(([entry]) => { onScreen = entry.isIntersecting; }).observe(terminal);

  // Every call to play() gets a new id; older, still-running calls see the
  // id change and stop, so clicking a dot interrupts the current example.
  let currentId = 0;

  const addLine = (html) => {
    const line = document.createElement("div");
    line.innerHTML = html || " ";
    screen.append(line);
    return line;
  };

  // Pauses while the terminal is scrolled away or the tab is in the background.
  const pause = async (ms, id) => {
    await sleep(ms);
    while ((!onScreen || document.hidden) && id === currentId) await sleep(250);
    return id === currentId;
  };

  async function play(index) {
    const id = ++currentId;
    const example = EXAMPLES[index];
    title.textContent = example.title;
    dots.forEach((dot, i) => dot.setAttribute("aria-current", i === index));
    screen.innerHTML = "";

    if (reduceMotion) {
      addLine(PROMPT + example.command);
      example.output.forEach((html) => addLine(html));
      return;
    }

    const typed = document.createElement("span");
    const cursor = document.createElement("span");
    cursor.className = "cursor";
    addLine(PROMPT).append(typed, cursor);

    for (const char of example.command) {
      if (!(await pause(35 + Math.random() * 55, id))) return;
      typed.textContent += char;
    }
    if (!(await pause(400, id))) return;
    cursor.remove();

    for (const html of example.output) {
      addLine(html);
      if (!(await pause(html.includes("›") ? 550 : 170, id))) return;
    }
    addLine("");
    addLine(PROMPT + '<span class="cursor"></span>');

    if (await pause(3500, id)) play((index + 1) % EXAMPLES.length);
  }

  play(0);
}

// ---------------------------------------------------------------- documentation page

function setupDocs() {
  const doc = $(".doc");
  if (!doc) return;
  const links = $$(".toc a");

  for (const heading of $$("h2[id], h3[id]", doc)) {
    heading.insertAdjacentHTML("beforeend", ` <a class="hash" href="#${heading.id}" aria-label="Link to this section">#</a>`);
  }

  // Highlight the sidebar link of the last section that scrolled past the top.
  const targets = links.map((link) => document.getElementById(link.hash.slice(1)));
  const highlightCurrent = () => {
    let current = links[0];
    targets.forEach((target, i) => {
      if (!target.closest("[hidden]") && target.getBoundingClientRect().top < 120) current = links[i];
    });
    links.forEach((link) => link.classList.toggle("active", link === current));
  };
  addEventListener("scroll", highlightCurrent, { passive: true });
  highlightCurrent();

  // Phones get a drop-down menu built from the same links.
  const select = $(".mobile-toc select");
  for (const group of $$(".toc-group")) {
    const optgroup = document.createElement("optgroup");
    optgroup.label = $("p", group).textContent;
    for (const link of $$("a", group)) optgroup.append(new Option(link.textContent, link.hash));
    select.append(optgroup);
  }
  select.addEventListener("change", () => {
    location.hash = select.value;
    select.selectedIndex = 0;
  });

  // Search hides every topic that does not mention the query.
  const input = $(".search input");
  const topics = $$("[data-topic]", doc);
  input.addEventListener("input", () => {
    const query = input.value.trim().toLowerCase();
    for (const topic of topics) {
      topic.hidden = query !== "" && !topic.textContent.toLowerCase().includes(query);
    }
    for (const section of $$(".doc-section", doc)) {
      const inner = $$("[data-topic]", section);
      if (inner.length) section.hidden = inner.every((t) => t.hidden);
    }
    for (const link of links) {
      const target = document.getElementById(link.hash.slice(1));
      link.parentElement.hidden = Boolean(target?.closest("[hidden]"));
    }
    for (const group of $$(".toc-group")) {
      group.hidden = $$("li", group).every((li) => li.hidden);
    }
    $(".no-results").hidden = topics.some((t) => !t.hidden);
  });

  addEventListener("keydown", (event) => {
    if (event.key === "/" && document.activeElement !== input) {
      event.preventDefault();
      input.focus();
    }
  });
}

setupNav();
setupReveal();
setupCounters();
setupCopyButtons();
setupTabs();
setupDownloads();
setupTypingTerminal();
setupDocs();
$$("[data-year]").forEach((el) => { el.textContent = new Date().getFullYear(); });
