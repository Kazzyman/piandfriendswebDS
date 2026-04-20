# Pi & Friends Web Suite

A collection of π calculation algorithms implemented in Go,
served as a live interactive web application.

Built and written entirely by **Richard (Rick) H. Woolley**.

## Live Demo

https://piandfriends.onrender.com

## What it does

Each method in the suite calculates π (or a related constant)
using a different algorithm, ranging from ancient geometry to
modern number theory. Also included is a brute-force method
for estimating the square or cube root of any integer to high
precision. The back-end is written in Go and streams results 
to the browser in real time using Server-Sent Events.

## Methods included

- **Archimedes** (c. 230 BCE) — polygon bisection, Rick's
  personal favorite. Produces 3,012 verified digits.
- **Spigot** — Rabinowitz-Wagon algorithm, integer arithmetic
  only, no floating point. Two-run honest edition with live
  uncertainty display and the Feynman Point easter egg.
- **BBP** — Bailey-Borwein-Plouffe, 1995. Calculates any
  individual hexadecimal digit of π without knowing the
  preceding digits.
- **Monte Carlo** — pure randomness converging on π. Rick's
  second favorite.
- **Chudnovsky** — the algorithm behind every world record.
  Independently verified to over 1,000 digits.
- **Custom Series** — a rapid alternating series.
- **Gauss-Legendre** — quadratic convergence. Correct digits
  double each iteration. Up to 4,930 verified digits.
- **Gregory-Leibniz** — 4 billion iterations, 10 digits.
- **Nilakantha** — c. 1530, predates Newton by 150 years.
  Two-phase: float64 then big.Float, breaking the float64
  wall live on screen.
- **Wallis** — 40 billion iterations, 10 digits. John Wallis,
  1655.
- **Brute-force Roots** — square and cube roots via rational
  approximation of perfect power pairs, inspired by ancient
  Greek methods. Includes the Delian Problem easter egg.
- **Euler's Number** — the natural logarithmic base e.

## Tech stack

- Back-end: Go (golang)
- Front-end: vanilla HTML, CSS, JavaScript
- Streaming: Server-Sent Events (SSE)
- Deployment: render.com

## Running locally

```bash
go run .
```

Then open http://localhost:8080 in your browser.

## Author

Richard (Rick) H. Woolley
https://github.com/Kazzyman
