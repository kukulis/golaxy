export class ShipDrawParams {
    constructor(drawMass, drawSpeed, drawGuns, drawAttack, drawDefence) {
        this.drawMass = drawMass;
        this.drawSpeed = drawSpeed;
        this.drawGuns = drawGuns;
        this.drawAttack = drawAttack;
        this.drawDefence = drawDefence;
    }
}

export class ShipDraw {
    static SVG_NS = "http://www.w3.org/2000/svg";

    static createSvgElement(tag, attrs) {
        const el = document.createElementNS(ShipDraw.SVG_NS, tag);
        for (const [key, value] of Object.entries(attrs)) {
            el.setAttribute(key, value);
        }
        return el;
    }


    scale = 10;
    color = "#666";
    circleLineColor = "#888";

    constructor(scale, color, circleLineColor) {
        this.scale = scale;
        this.color = color;
        this.circleLineColor = circleLineColor;
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

    drawShipRaw(svg, cx, cy, drawMass, drawSpeed, drawGuns, drawAttack, drawDefence, rotation = 0) {
        const s = this.scale;
        const scaledCircleRadius = drawMass * s;
        const scaledTrapezeHeight = drawSpeed * s;
        const scaledCombSpacing = s*2;
        const scaledCombTeethLength = drawAttack * s;

        const circleCx = cx + scaledTrapezeHeight + scaledCircleRadius;
        const trapezeCx = circleCx - scaledCircleRadius - scaledTrapezeHeight / 2;
        const combCx = circleCx + scaledCircleRadius;

        // Create wrapper group with rotation around input position (cx, cy)
        const group = ShipDraw.createSvgElement("g", {
            transform: `rotate(${rotation}, ${cx}, ${cy})`
        });
        // Alternative: rotate around circle center
        // const group = ShipDraw.createSvgElement("g", {
        //     transform: `rotate(${rotation}, ${circleCx}, ${cy})`
        // });

        this.drawTrapeze(group, trapezeCx, cy, scaledTrapezeHeight, this.color);
        this.drawCircle(group, circleCx, cy, scaledCircleRadius, this.color, this.circleLineColor, drawDefence*this.scale/5);
        this.drawComb(group, combCx, cy, scaledCombSpacing, drawGuns, scaledCombTeethLength, this.color);

        svg.appendChild(group);
    }
}
