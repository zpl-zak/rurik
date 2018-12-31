# Rurik Framework

A 2D cross-platform game engine/framework made as an experiment to explore what possibilities can be achieved in an ideal workspace.

This framework is in active development and its features/API is subject to change. Any PRs are welcomed, especially those improving code quality and performance.

## Features

Rurik Frammework consists of the following key features:
- Fast, small and clean codebase, which is easily extendable and tweakable.
- [Tiled](https://www.mapeditor.org/) map support which makes level designing much easier.
- Robust scripting backend, which offers an ability to flexibly alternate the game world or introduce scripted sequences.
- Adaptive save system with ability to extend it's support onto custom game data which needs to be kept persistent.
- Music manager for your musical needs.
- Basic AABB collision detection and resolution.
- Fast frustum-culled renderer offering great performance under heavier loads.
- Various built-in entity classes encapsulating stereotypes, such as trigger zones, collision areas, timers or even dialogue emitters.
- Straightforward dialogue system.
- Basic lightmap generator, currently supporting additive and multiplicative lighting solutions.
- Weather system, supports 4 time of day stages and influences the whole game environment.
- Simple set of tools to profile parts of your game logic and display custom statistics in an editor UI.
- Currently runs on Linux, Windows and macOS.
- Ability to easily render to texture or manipulate your render target (blur, ...).

## Future plans

While the framework already offers interesting features, I plan to expand it's features with the following ideas:
- Add support for Android, iOS and WebAssembly.
- Expand the weather system, implement various weather effects.
- Solidify the filesystem API, to make it easier for us when adding Android support.
- Revamp the entity system, add components.
- Refactor some entity classes into actual core features.
- Make better use of the scripting backend.
- Bring my own assets and avoid using third-party assets for the demo.
- Add sound support + sound mixing.
- Improve the demo code to showcase framework's features.

![](https://user-images.githubusercontent.com/9026786/50441112-738cd880-08f9-11e9-95fd-4e0d074bcb20.png)
![](https://i.imgur.com/6b98kOA.png)

## How to build

Navigate to `src/demo` and execute `go get ./...` to fetch all dependencies. Afterwards, navigate back to the root folder and execute `make` to build the game.

Make sure you meet all the requirements at [raylib-go](https://github.com/zaklaus/raylib-go) before you compile the project.

## License

See [COPYING.md](COPYING.md) for licensing information. The demo code residing in `src/demo` directory falls under public domain and can be used freely.

Used assets are licensed under their respective licenses. Unless explicitly stated, a written permission of its authors is required for the content to be used outside of this repository.
