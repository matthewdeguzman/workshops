# Go Workshop - Awesome QR Code generator

This workshop goes over creating the backeend for a QR code generator using Go.

# Demo

https://github.com/user-attachments/assets/68a6e9b3-65cc-43cf-b525-6d227d189e40

# Running the project

__Prerequisites__:
- NodeJS - Follow [svelte's guide](https://svelte.dev/blog/svelte-for-new-developers#installing-node-js) to installing NodeJS
- [Install Go](https://go.dev/doc/install)

## Backend
1. Navigate to the backend directory and run `go cmd/app/main.go` (optionally, if you have `make` installed run `make run`)

## Frontend

1. Navigate to the frontend directory and run `npm i`
2. Create a file `.env` in the frontend directory and add the line `VITE_BASE_URL=http://localhost:8080`
3. Run `npm run dev`

The front-end should successfully start and be accessible through `http://localhost:5173`.
