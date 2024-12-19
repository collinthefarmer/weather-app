initObservationCanvas();

function initObservationCanvas() {
    const customElementRegistry = window.customElements;
    if (customElementRegistry === null) {
        alert("no custom elements? :(");
        return;
    }

    class ObservationCanvas extends HTMLElement {
        constructor() {
            super();
        }

        connectedCallback() {
            this.style = "position: relative";

            const shadow = this.attachShadow({ mode: "open" });
            const canvas = document.createElement("canvas");
            canvas.style = "position: absolute; inset: 0 0 0 0;";
            canvas.width = 500;
            canvas.height = 500;
            shadow.appendChild(canvas);

            this.canvas = canvas;
            this.cursor = new ObservationCanvasBrush(canvas, 10);

            canvas.addEventListener("mousedown", () => {
                this.cursor.start();
            });

            canvas.addEventListener("mousemove", (ev) => {
                this.cursor.paint(
                    ev.clientX - this.offsetLeft,
                    ev.clientY - this.offsetTop,
                );
            });

            canvas.addEventListener("mouseup", () => {
                this.cursor.stop();
            });
        }
    }

    class ObservationCanvasBrush {
        constructor(canvas, size, defaultColor = "#000000") {
            this.canvas = canvas;
            this.context = canvas.getContext("2d");

            const paintBufferSize = 255;
            this.paintBuffer = new Uint16Array(paintBufferSize);
            this.paintBufferPos = 0;

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

        applyPaintBuffer() {
            this.context.ellipse(x, y, this.size, this.size, 0, 0, 2 * Math.PI);
            this.context.fill();
        }
    }

    class ObservationCanvasPallete extends HTMLElement {
        static observedAttributes = ["for"];

        constructor() {
            super();
        }

        attributeChangedCallback(attr, oldval, newval) {
            switch (attr) {
                case "for":
                    if (oldval) {
                        this.disconnectCanvas(oldval);
                    }
                    this.connectCanvas(newval);
                    break;
            }
        }

        disconnectCanvas(id) {
            const canvas = document.querySelector("#" + id);
            console.log(canvas);
        }

        connectCanvas(id) {
            const canvas = document.querySelector("#" + id);
            console.log(canvas);
        }
    }

    customElementRegistry.define("observation-canvas", ObservationCanvas);
    customElementRegistry.define(
        "observation-canvas-pallete",
        ObservationCanvasPallete,
    );
}
