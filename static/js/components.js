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
    }

    customElementRegistry.define("observation-canvas", ObservationCanvas);
}
