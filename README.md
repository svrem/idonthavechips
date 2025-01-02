# IDontHaveChips

IDontHaveChips is a program written in Go and HTML that allows users to create and play in online betting session, where they can bet with fake money. The host creates a session and the players can then join that session using the Session ID presented on the host's screen. The host can then start the game and the players can place their bets. The winners are determined by the host and the players can then see their winnings.

## Installation

### Requirements

- [Go](https://golang.org/dl/)
- [Git](https://git-scm.com/downloads)
- [Node.js](https://nodejs.org/en/download/) (dev)
- [NPM](https://www.npmjs.com/get-npm) (dev)

### Steps

1. Clone the repository
```bash
git clone https://github.com/svrem/IDontHaveChips.git && cd IDontHaveChips
```

2. Install the dependencies
```bash
npm install
```

3. Run the server
```bash
go run .
```

4. Open the browser and go to `localhost:8080`

### Alternative (Docker)

1. Build the Docker image
```bash
docker build -t idonthavechips .
```

2. Run the Docker container
```bash
docker run -p 8080:8080 idonthavechips
```

3. Open the browser and go to `localhost:8080`


