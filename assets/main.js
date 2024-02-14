// nostalgia for a simpler, more complicated time
const $ = document.querySelector.bind(document);
const $$ = document.querySelectorAll.bind(document);

const initialLines = 10;

class Button extends HTMLButtonElement {
  constructor() {
    super();
    this.addEventListener("click", this.click);
  }
}

class LineRemover extends Button {
  click() {
    this.closest("div.linecontainer").parentElement.remove();
  }
}

class LinePinner extends Button {
  click() {
    const l = this.closest("div.linecontainer");
    l.classList.toggle("unpinned");
    if (l.classList.contains("unpinned")) {
      this.innerText = "pin";
    } else {
      this.innerText = "upin";
    }
  }
}

class LineUpper extends Button {
  click() {
    // TODO debug why when I:
    // - add a new line
    // - move that line up
    // the "down" button stays disabled
    const l = this.closest("div.linecontainer").parentElement;
    const s = l.previousElementSibling;
    s.before(l);
    this.checkDisabled();
    s.querySelector("button[is=line-upper]").checkDisabled()
    s.querySelector("button[is=line-downer]").checkDisabled()
  }

  checkDisabled() {
    const l = this.closest("div.linecontainer").parentElement;
    if (l.previousElementSibling == null) {
      this.setAttribute("disabled", "yeah");
    } else {
      this.removeAttribute("disabled");
    }
  }

  connectedCallback() {
    this.checkDisabled();
  }
}

class LineDowner extends Button {
  click() {
    const l = this.closest("div.linecontainer").parentElement;
    const s = l.nextElementSibling;
    s.after(l);
    this.checkDisabled()
    s.querySelector("button[is=line-downer]").checkDisabled()
    s.querySelector("button[is=line-upper]").checkDisabled()
  }

  checkDisabled() {
    const l = this.closest("div.linecontainer").parentElement;
    if (l.nextElementSibling == null) {
      this.setAttribute("disabled", "yeah");
    } else {
      this.removeAttribute("disabled");
    }
  }

  connectedCallback() {
    if (this.connected) {
      return;
    }
    const count = $$("button[is=line-downer]").length;
    if (count != initialLines) {
      return;
    }
    this.checkDisabled();
    this.connected = true;
  }
}

class LineAdder extends Button {
  click() {
    $("div[is=lines-div]").add()
  }
}

class PoemRegenner extends Button {
  click() {
    $$(".unpinned").forEach((e) => {
      e.parentElement.regen();
    });
  }
}

class PoemLine extends HTMLDivElement {
  constructor() {
    super();
    this.ltp = $("#linetmpl");
  }

  connectedCallback() {
    if (this.connected) {
      return;
    }
    this.appendChild(this.ltp.content.cloneNode(true));
    this.connected = true;
  }

  regen() {
    fetch(new Request("/line")).then((resp) => {
      if (!resp.ok) {
        throw new Error(`sadness stalks the land in ${resp.status} ways`);
      }
      return resp.json();
    }).then((payload) => {
      this.querySelector(".linetext").innerText = payload.Text;
      this.querySelector(".linetext").setAttribute("data-source", payload.Source.Name);
    });
  }
}

class Lines extends HTMLDivElement {
  constructor() {
    super();
  }

  connectedCallback() {
    for (var i = 0; i < initialLines; i++) {
      this.add();
    }
    addEventListener("beforeunload", (e) => {
      if ($$("div.linecontainer:not(.unpinned)").length > 0) {
        e.preventDefault();
      }
    });
  }

  add() {
    var ld = document.createElement("div", {is: "poem-line"});
    this.append(ld);
    ld.regen();
    this.querySelectorAll("button[is=line-downer]").forEach((e) => {
      e.checkDisabled();
    })
  }
}

customElements.define("line-remover", LineRemover, { extends: "button" });
customElements.define("line-pinner", LinePinner, { extends: "button" });
customElements.define("line-upper", LineUpper, { extends: "button" });
customElements.define("line-downer", LineDowner, { extends: "button" });
customElements.define("line-adder", LineAdder, { extends: "button" });
customElements.define("poem-regenner", PoemRegenner, {extends: "button"});
customElements.define("poem-line", PoemLine, {extends: "div"});
customElements.define("lines-div", Lines, {extends: "div"});
