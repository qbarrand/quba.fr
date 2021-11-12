import {EventEmitter} from 'events';
import {Constraint, Orientation} from "./constraint";

const p = new Promise<boolean>((resolve, _) => {
    const img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = 'data:image/webp;base64,UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA'
});

export class ImageChangedEvent {
    constructor(readonly url: string, readonly date: string, readonly location: string, readonly mainColor: string) {}
}

export class BackgroundManager {
    private currentConstraint: Constraint
    private readonly emitter = new EventEmitter()

    constructor(private name: string) {}

    addEventListener(name: string, fn) {
        this.emitter.on(name, fn)
    }

    async updateConstraint(c: Constraint) {
        console.debug(`trying to update constraint with ${c}`)

        if (!this.requiresUpdate(c)) {
            return
        }

        const params = {format: await p ? 'webp':'jpg',}

        if (this.name != null) {
            params['name'] = this.name
        }

        if (c.high != Infinity) {
            switch (c.orientation) {
                case Orientation.Landscape:
                    params['width'] = c.high
                    break
                case Orientation.Portrait:
                    params['height'] = c.high
                    break
            }
        }

        const f = await fetch('/background?' + new URLSearchParams(params))
        const json = await f.json()

        this.name = json.name

        const event = new ImageChangedEvent(
            json.url,
            json.date,
            json.location,
            json.mainColor
        )

        this.emitter.emit('change', event)

        this.currentConstraint = c
    }

    private requiresUpdate(c: Constraint): boolean {
        return this.currentConstraint === undefined ||
            this.currentConstraint.orientation != c.orientation ||
            this.currentConstraint.high < c.high
    }
}