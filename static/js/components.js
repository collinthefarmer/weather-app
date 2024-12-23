initObservationCanvas();

function initObservationCanvas() {
    const customElementRegistry = window.customElements;
    if (customElementRegistry === null) {
        alert("no custom elements? :(");
        return;
    }

    class ObservationCanvas extends HTMLElement {
        static observedAttributes = ["width", "height"];

        pallete;

        constructor() {
            super();
            this.attachShadow({ mode: "open" });

            this.width = 0;
            this.height = 0;

            this.canvas = document.createElement("canvas");
            this.canvas.id = "canvas";

            this.brush = new ObservationCanvasBrush(this, 10);

            this.dataBuffer = new Uint8Array();
        }

        connectedCallback() {
            this.canvas.addEventListener("mousedown", () => {
                this.brush.start();
            });

            this.canvas.addEventListener("mousemove", (ev) => {
                if (this.brush.painting) {
                    const x = ev.clientX - this.offsetLeft;
                    const y = ev.clientY - this.offsetTop;
                    this.brush.paint(x, y);
                    this.#addData(x, y, this.pallete.color, this.pallete.size);
                }
            });

            this.canvas.addEventListener("mouseup", () => {
                if (this.brush.painting) {
                    this.brush.stop();
                    this.#serializeData();
                }
            });

            this.canvas.addEventListener("mouseleave", () => {
                if (this.brush.painting) {
                    this.brush.stop();
                    this.#serializeData();
                }
            });

            this.shadowRoot.appendChild(this.canvas);
        }

        attributeChangedCallback(attr, _oldval, newval) {
            switch (attr) {
                case "height":
                    this.height = parseInt(newval.replace(/[^\d]/g, ""), 10);
                    this.canvas.height = this.height;
                    this.#resetData();
                    break;
                case "width":
                    this.width = parseInt(newval.replace(/[^\d]/g, ""), 10);
                    this.canvas.width = this.width;
                    this.#resetData();
                    break;
            }
        }

        #resetData() {
            this.dataBuffer = new Uint8Array(this.height * this.width);
        }

        #addData(x, y, iColor, iSize) {
            this.dataBuffer.set([iColor % 16, iSize % 16], x + y * this.width);
        }

        #serializeData() {}
    }

    class ObservationCanvasBrush {
        constructor(canvas) {
            this.canvas = canvas;
            this.context = canvas.canvas.getContext("2d");

            this.painting = false;
        }

        start() {
            this.painting = true;
            this.context.beginPath();

            this.context.fillStyle =
                this.canvas.pallete.colors[this.canvas.pallete.color];
        }

        paint(x, y) {
            if (this.painting) {
                this.context.ellipse(
                    x,
                    y,
                    this.canvas.pallete.sizes[this.canvas.pallete.size],
                    this.canvas.pallete.sizes[this.canvas.pallete.size],
                    0,
                    0,
                    2 * Math.PI,
                );
                this.context.fill();
            }
            this.context.moveTo(x, y);
        }

        stop() {
            this.painting = false;
        }
    }

    class ObservationCanvasPallete extends HTMLElement {
        static observedAttributes = ["for", "colorset", "brushset"];

        canvas;

        colors;
        color = 0;

        sizes;
        size = 0;

        constructor() {
            super();
            this.attachShadow({ mode: "open" });
        }

        connectedCallback() {
            const sheet = new CSSStyleSheet();

            sheet.insertRule(
                ":host { display: flex; flex-flow: column; align-items: stretch; gap: .5rem; }",
            );
            sheet.insertRule(
                "fieldset { display: flex; flex-flow: row; justify-content: stretch; align-items: stretch; padding: 0; }",
            );
            sheet.insertRule(
                "fieldset input { padding: 0; margin: 0; appearance: none; aspect-ratio: 1; flex: 1; position: relative; }",
            );
            sheet.insertRule(
                "fieldset input:checked { outline: solid; z-index: 2; }",
            );

            sheet.insertRule(
                "#brushset input::after { content: attr(value); position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; }",
            );

            this.shadowRoot.adoptedStyleSheets = [sheet];
        }

        attributeChangedCallback(attr, _oldval, newval) {
            switch (attr) {
                case "for":
                    this.canvas = document.querySelector("#" + newval);
                    this.canvas.pallete = this;
                    break;
                case "colorset":
                    this.applyColorset(newval ?? "");
                    break;
                case "brushset":
                    this.applyBrushset(newval ?? "");
                    break;
            }
        }

        applyColorset(colorsetString) {
            this.colors = colorsetString.split(/, ?/);

            let fieldset = this.shadowRoot.querySelector("#colorset");
            if (!fieldset) {
                fieldset = document.createElement("fieldset");
                fieldset.id = "colorset";

                this.shadowRoot.appendChild(fieldset);
            } else {
                fieldset.innerHTML = "";
            }

            for (let i = 0; i < this.colors.length; i++) {
                const colorInput = document.createElement("input");
                colorInput.type = "radio";
                colorInput.name = "paint-color";
                colorInput.value = this.colors[i];
                colorInput.style.setProperty(
                    "background-color",
                    this.colors[i],
                );

                const changeBrushColor = () => (this.color = i);
                colorInput.addEventListener("click", changeBrushColor);
                if (i === 0) {
                    changeBrushColor();
                    colorInput.checked = true;
                }

                fieldset.appendChild(colorInput);
            }
        }

        applyBrushset(brushsetString) {
            this.sizes = brushsetString.split(/, ?/);

            let fieldset = this.shadowRoot.querySelector("#brushset");
            if (!fieldset) {
                fieldset = document.createElement("fieldset");
                fieldset.id = "brushset";

                this.shadowRoot.appendChild(fieldset);
            } else {
                fieldset.innerHTML = "";
            }

            for (let i = 0; i < this.sizes.length; i++) {
                const brushInput = document.createElement("input");
                brushInput.type = "radio";
                brushInput.name = "brush-size";
                brushInput.value = this.sizes[i];

                const changeBrushSize = () => (this.size = i);
                brushInput.addEventListener("click", changeBrushSize);
                if (i === 0) {
                    changeBrushSize();
                    brushInput.checked = true;
                }

                fieldset.appendChild(brushInput);
            }

            this.shadowRoot.appendChild(fieldset);
        }
    }

    customElementRegistry.define("observation-canvas", ObservationCanvas);
    customElementRegistry.define(
        "observation-canvas-pallete",
        ObservationCanvasPallete,
    );
}
