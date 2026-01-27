class ShipDraw {
    static SVG_NS = "http://www.w3.org/2000/svg";

    static createSvgElement(tag, attrs) {
        const el = document.createElementNS(ShipDraw.SVG_NS, tag);
        for (const [key, value] of Object.entries(attrs)) {
            el.setAttribute(key, value);
        }
        return el;
    }


    scale = 10;

    constructor(scale) {
        this.scale = scale
    }

    drawTrapezeRaw(svg, cx, cy, height, topWidth, bottomWidth, color) {
        const halfHeight = height / 2;
        const halfTop = topWidth / 2;
        const halfBottom = bottomWidth / 2;

        const points = [
            `${cx - halfHeight},${cy - halfTop}`,
            `${cx + halfHeight},${cy - halfBottom}`,
            `${cx + halfHeight},${cy + halfBottom}`,
            `${cx - halfHeight},${cy + halfTop}`
        ].join(" ");

        const polygon = ShipDraw.createSvgElement("polygon", {
            points: points,
            fill: color,
            stroke: color,
            "stroke-width": "2"
        });
        svg.appendChild(polygon);
    }

    drawTrapeze(svg, cx, cy, height, color) {
        const topWidth = height * 0.5;
        const bottomWidth = height * 1.5
        this.drawTrapezeRaw(svg, cx, cy, height, bottomWidth, topWidth, color);
    }

    drawCircle(svg, cx, cy, radius, fillColor, strokeColor, strokeWidth) {
        const circle = ShipDraw.createSvgElement("circle", {
            cx: cx,
            cy: cy,
            r: radius,
            fill: fillColor,
            stroke: strokeColor,
            "stroke-width": strokeWidth
        });
        svg.appendChild(circle);
    }

    drawComb(svg, cx, cy, spacing, teethCount, teethLength, color) {
        const group = ShipDraw.createSvgElement("g", {
            stroke: color,
            "stroke-width": "3",
            "stroke-linecap": "round"
        });

        const baseHeight = teethCount > 1 ? spacing * (teethCount - 1) : 0;
        const halfBase = baseHeight / 2;
        const baseLine = ShipDraw.createSvgElement("line", {
            x1: cx,
            y1: cy - halfBase,
            x2: cx,
            y2: cy + halfBase
        });
        group.appendChild(baseLine);

        for (let i = 0; i < teethCount; i++) {
            const y = cy - halfBase + i * spacing;
            const tooth = ShipDraw.createSvgElement("line", {
                x1: cx,
                y1: y,
                x2: cx + teethLength,
                y2: y
            });
            group.appendChild(tooth);
        }

        svg.appendChild(group);
    }

    drawCombined(svg, cx, cy, circleRadius, trapezeHeight, combTeethCount, combTeethLength, color, circleLineColor) {
        const s = this.scale;
        const scaledCircleRadius = circleRadius * s;
        const scaledTrapezeHeight = trapezeHeight * s;
        const scaledCombSpacing = s;
        const scaledCombTeethLength = combTeethLength * s;

        const circleCx = cx + scaledTrapezeHeight + scaledCircleRadius;
        const trapezeCx = circleCx - scaledCircleRadius - scaledTrapezeHeight / 2;
        const combCx = circleCx + scaledCircleRadius;

        this.drawTrapeze(svg, trapezeCx, cy, scaledTrapezeHeight, color);
        this.drawCircle(svg, circleCx, cy, scaledCircleRadius, color, circleLineColor, 2);
        this.drawComb(svg, combCx, cy, scaledCombSpacing, combTeethCount, scaledCombTeethLength, color);
    }
}
