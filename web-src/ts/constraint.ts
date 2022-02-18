export class Orientation {
    readonly name: string
    readonly mqdim: string

    constructor(name: string, mqdim: string) {
        this.name = name
        this.mqdim = mqdim
    }
}

export const LANDSCAPE = new Orientation('landscape', 'width')
export const PORTRAIT = new Orientation('portrait', 'height')

export class Constraint {
    readonly orientation: Orientation
    readonly low: number
    readonly high: number

    constructor(o: Orientation, l: number, h: number) {
        this.orientation = o
        this.low = l
        this.high = h
    }

    toMediaQuery(): string {
        let mq = `(orientation: ${this.orientation.name})`

        if (this.low != 0) {
            mq += ` and (min-${this.orientation.mqdim}: ${this.low}px)`
        }

        if (this.high != Infinity) {
            mq += ` and (max-${this.orientation.mqdim }: ${this.high}px)`
        }

        return mq
    }
}

type MediaConstraint = {
    [key: string]: Constraint
}

function getConstraints(o: Orientation, n: number[]): MediaConstraint {
    const mcs: MediaConstraint = {}

    let previous = 0

    for (let i = 0; i < n.length; i++) {
        if (i != 0) {
            previous++
        }

        const v = n[i]
        const c = new Constraint(o, previous, v)

        mcs[c.toMediaQuery()] = c

        previous = v
    }

    if (n.length > 0) {
        const c = new Constraint(o, previous+1, Infinity)
        mcs[c.toMediaQuery()] = c
    }

    return mcs
}

export function generateMediaConstraints(widths: number[], heights: number[]): MediaConstraint {
    return {...getConstraints(LANDSCAPE, widths), ...getConstraints(PORTRAIT, heights)}
}
