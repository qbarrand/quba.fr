import {EventEmitter} from "events";

export enum Orientation {
    Landscape,
    Portrait
}

export class Constraint {
    private readonly emitter = new EventEmitter()
    readonly orientation: Orientation
    readonly low: number
    readonly high: number

    constructor(o: Orientation, l: number, h: number) {
        this.orientation = o
        this.low = l
        this.high = h

        window.addEventListener('resize', this.processResize.bind(this))
    }

    addEventListener(name: string, fn) {
        this.emitter.addListener(name, fn)
    }

    processResize() {
        console.debug(`onresize triggered for ${this}`)

        const scale = window.devicePixelRatio

        const w = Math.floor(innerWidth * scale)
        const h = Math.floor(innerHeight * scale)

        console.debug(`innerWidth: ${innerWidth}, innerHeight: ${innerHeight}`)
        console.debug(`w: ${w}, h: ${h}, scale: ${scale}`)

        const or = w < h ? Orientation.Portrait : Orientation.Landscape

        if (or != this.orientation) {
            return
        }

        if (
            or == this.orientation && (
                (this.orientation == Orientation.Landscape && this.low <= w && w <= this.high) ||
                (this.orientation == Orientation.Portrait && this.low <= h && h <= this.high)
            )
        ) {
            this.emitter.emit('active', this)
        }
    }

    toString(): string {
        return `${Orientation[this.orientation]}, ${this.low} -- ${this.high}`
    }
}

function getConstraints(o: Orientation, n: number[]): Constraint[] {
    const mcs: Constraint[] = []

    let previous = 0

    for (let i = 0; i < n.length; i++) {
        if (i != 0) {
            previous++
        }

        const v = n[i]

        mcs.push(
            new Constraint(o, previous, v)
        )

        previous = v
    }

    if (n.length > 0) {
        mcs.push(
            new Constraint(o, previous+1, Infinity)
        )
    }

    return mcs
}

export function generateMediaConstraints(widths: number[], heights: number[]): Constraint[] {
    return [
        ...getConstraints(Orientation.Landscape, widths),
        ...getConstraints(Orientation.Portrait, heights)
    ]
}
