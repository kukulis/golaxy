const SVG_NS = "http://www.w3.org/2000/svg";

function createSvgElement(tag, attrs) {
    const el = document.createElementNS(SVG_NS, tag);
    for (const [key, value] of Object.entries(attrs)) {
        el.setAttribute(key, value);
    }
    return el;
}

function drawTrapezeRaw(svg, cx, cy, height, topWidth, bottomWidth, color) {
    const halfHeight = height / 2;
    const halfTop = topWidth / 2;
    const halfBottom = bottomWidth / 2;

    const points = [
        `${cx - halfHeight},${cy - halfTop}`,
        `${cx + halfHeight},${cy - halfBottom}`,
        `${cx + halfHeight},${cy + halfBottom}`,
        `${cx - halfHeight},${cy + halfTop}`
    ].join(" ");

    const polygon = createSvgElement("polygon", {
        points: points,
        fill: color,
        stroke: color,
        "stroke-width": "2"
    });
    svg.appendChild(polygon);
}

function drawTrapeze(svg, cx, cy, height, color) {
    const topWidth = height * 0.5;
    const bottomWidth = height * 1.5
    drawTrapezeRaw(svg, cx, cy, height, bottomWidth, topWidth, color);
}

function drawCircle(svg, cx, cy, radius, fillColor, strokeColor, strokeWidth) {
    const circle = createSvgElement("circle", {
        cx: cx,
        cy: cy,
        r: radius,
        fill: fillColor,
        stroke: strokeColor,
        "stroke-width": strokeWidth
    });
    svg.appendChild(circle);
}

function drawComb(svg, cx, cy, spacing, teethCount, teethLength, color) {
    const group = createSvgElement("g", {
        stroke: color,
        "stroke-width": "3",
        "stroke-linecap": "round"
    });

    const baseHeight = teethCount > 1 ? spacing * (teethCount - 1) : 0;
    const halfBase = baseHeight / 2;
    const baseLine = createSvgElement("line", {
        x1: cx,
        y1: cy - halfBase,
        x2: cx,
        y2: cy + halfBase
    });
    group.appendChild(baseLine);

    for (let i = 0; i < teethCount; i++) {
        const y = cy - halfBase + i * spacing;
        const tooth = createSvgElement("line", {
            x1: cx,
            y1: y,
            x2: cx + teethLength,
            y2: y
        });
        group.appendChild(tooth);
    }

    svg.appendChild(group);
}

function drawCombined(svg, cy, circleRadius, trapezeHeight, combSpacing, combTeethCount, combTeethLength, color, circleLineColor) {
    const circleCx = trapezeHeight + circleRadius;
    const trapezeCx = circleCx - circleRadius - trapezeHeight / 2;
    const combCx = circleCx + circleRadius;

    drawTrapeze(svg, trapezeCx, cy, trapezeHeight, color);
    drawCircle(svg, circleCx, cy, circleRadius, color, circleLineColor, 2);
    drawComb(svg, combCx, cy, combSpacing, combTeethCount, combTeethLength, color);
}
