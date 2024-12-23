initObservationCanvas();

function initObservationCanvas() {
    const customElementRegistry = window.customElements;
    if (customElementRegistry === null) {
        alert("no custom elements? :(");
        return;
    }

    class ObservationCanvas extends HTMLElement {
        static observedAttributes = ["width", "height"];

        constructor() {
            super();
            this.attachShadow({ mode: "open" });

            this.width = 0;
            this.height = 0;

            this.canvas = document.createElement("canvas");
            this.canvas.id = "canvas";

            this.brush = new ObservationCanvasBrush(this.canvas, 10);
        }

        connectedCallback() {
            this.canvas.addEventListener("mousedown", () => {
                this.brush.start();
            });

            this.canvas.addEventListener("mousemove", (ev) => {
                this.brush.paint(
                    ev.clientX - this.offsetLeft,
                    ev.clientY - this.offsetTop,
                );
            });

            this.canvas.addEventListener("mouseup", () => {
                this.brush.stop();
            });

            this.canvas.addEventListener("mouseleave", () => {
                this.brush.stop();
            });

            this.shadowRoot.appendChild(this.canvas);
        }

        attributeChangedCallback(attr, _oldval, newval) {
            switch (attr) {
                case "height":
                    this.height = parseInt(newval.replace(/[^\d]/g, ""), 10);
                    this.canvas.height = this.height;
                    break;
                case "width":
                    this.width = parseInt(newval.replace(/[^\d]/g, ""), 10);
                    this.canvas.width = this.width;
                    break;
            }
        }
    }

    class ObservationCanvasBrush {
        constructor(canvas, size, defaultColor = "#000000") {
            this.canvas = canvas;
            this.context = canvas.getContext("2d");

            this.painting = false;

            this.size = size;
            this.color = defaultColor;
        }

        start() {
            this.painting = true;
            this.context.beginPath();
        }

        paint(x, y) {
            if (this.painting) {
                this.context.ellipse(
                    x,
                    y,
                    this.size,
                    this.size,
                    0,
                    0,
                    2 * Math.PI,
                );
                this.context.fillStyle = this.color;
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
                    break;
                case "colorset":
                    this.applyColorset(newval ?? "");
                    break;
                case "brushset":
                    this.applyBrushset(newval ?? "");
                    break;
            }
        }

        canvasBrushChanger(property, value) {
            return () => {
                if (
                    this.canvas &&
                    this.canvas.nodeName === "OBSERVATION-CANVAS"
                ) {
                    this.canvas.brush[property] = value;
                } else {
                    console.error("pallete is not for a valid canvas");
                }
            };
        }

        applyColorset(colorsetString) {
            const colorset = colorsetString.split(/, ?/);

            let fieldset = this.shadowRoot.querySelector("#colorset");
            if (!fieldset) {
                fieldset = document.createElement("fieldset");
                fieldset.id = "colorset";

                this.shadowRoot.appendChild(fieldset);
            }

            for (let i = 0; i < colorset.length; i++) {
                const colorInput = document.createElement("input");
                colorInput.type = "radio";
                colorInput.name = "paint-color";
                colorInput.value = colorset[i];
                colorInput.style.setProperty("background-color", colorset[i]);

                const clickFunc = this.canvasBrushChanger("color", colorset[i]);

                colorInput.addEventListener("click", clickFunc);
                if (i === 0) {
                    clickFunc();
                    colorInput.checked = true;
                }

                fieldset.appendChild(colorInput);
            }
        }

        applyBrushset(brushsetString) {
            const brushset = brushsetString.split(/, ?/);

            let fieldset = this.shadowRoot.querySelector("#brushset");
            if (!fieldset) {
                fieldset = document.createElement("fieldset");
                fieldset.id = "brushset";

                this.shadowRoot.appendChild(fieldset);
            }

            for (let i = 0; i < brushset.length; i++) {
                const brushInput = document.createElement("input");
                brushInput.type = "radio";
                brushInput.name = "brush-size";
                brushInput.value = brushset[i];

                const clickFunc = this.canvasBrushChanger(
                    "size",
                    parseInt(brushset[i]),
                );
                brushInput.addEventListener("click", clickFunc);
                if (i === 0) {
                    clickFunc();
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
