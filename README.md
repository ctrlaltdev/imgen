# IMGEN

Image Generator

## Build to Docker

```sh
docker build . -t imgen
```

## Run Docker

```sh
docker run -d --name imgen -p 3000:3000 imgen
```

You can also customize what port is used in the container:
```sh
docker run -d --name imgen -e "PORT=5000" -p 3000:5000 imgen
```

## Usage

There are 3 parameters available:

- Format (svg, png, jpg), default: svg
- Width, default: 1920
- Height, default: 1080

You can specify only the format, or only the dimensions, both or none:

- `/`
- `/png/`
- `/300/150/`
- `/jpg/200/600/`
