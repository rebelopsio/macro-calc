# Macro Calculator

A web-based macro nutrient calculator using the Mifflin-St Jeor equation. Built with Go, Templ, HTMX, AlpineJS, and Tailwind CSS.

## Features

- Calculate daily macronutrients based on:
  - Sex, age, height, and weight
  - Activity level
  - Fitness goals (lose, maintain, or gain weight)
- Advanced weekly planning with different activity levels per day
- Dark/light mode toggle
- Responsive design for all devices
- Army green color palette

## Tech Stack

- **Backend**: Go with Echo framework
- **Templating**: Templ
- **Frontend**: HTMX + AlpineJS
- **Styling**: Tailwind CSS
- **Deployment**: Netlify

## Development

### Prerequisites

- Go 1.21+
- Node.js (for Tailwind CSS)
- Make

### Setup

1. Clone the repository:
```bash
git clone https://github.com/rebelopsio/macro-calc.git
cd macro-calc
```

2. Install dependencies:
```bash
make install-deps
```

3. Run in development mode:
```bash
make dev
```

In another terminal, run CSS watcher:
```bash
make css-watch
```

4. Open http://localhost:8080

### Building for Production

```bash
make build
```

### Running Tests

```bash
make test
```

## Deployment

This project is configured for deployment on Netlify. The `netlify.toml` file contains the necessary configuration.

## License

MIT