# Awesome Hytale

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The complete directory of Hytale community resources. Find mods, servers, tools, wikis, and more.

**Website:** [awesomehytale.com](https://www.awesomehytale.com)

## Categories

- [Official Channels](resources/official-channels/) - Official Hytale presence
- [Wikis & Knowledge](resources/wikis-knowledge/) - Wikis and guide sites
- [Forums & Discussion](resources/forums-discussion/) - Community forums and Discord
- [News & Media](resources/news-media/) - News sites and coverage
- [YouTube & Creators](resources/youtube-creators/) - Content creators
- [Tools & Utilities](resources/tools-utilities/) - Tools, APIs, and bots
- [Development & Modding](resources/development-modding/) - Modding resources
- [Mods & Assets](resources/mods-assets/) - Mod repositories
- [Servers](resources/servers/) - Server lists and hosting
- [Social & Communities](resources/social-communities/) - Fan communities
- [Creator Program](resources/creator-program/) - Creator program info
- [Other Resources](resources/other-resources/) - Miscellaneous resources

## Contributing

We welcome contributions! To add a new resource:

1. Fork this repository
2. Create a new markdown file in the appropriate category folder
3. Follow the resource template format below
4. Submit a pull request

### Resource Template

```markdown
# Resource Name

**Website:** [https://example.com](https://example.com)

**Category:** Category Name > Subcategory Name

---

## Overview

Description of the resource.

## Details

| Property | Value |
|----------|-------|
| **Platform** | Web/Discord/etc. |
| **Audience** | All/Modders/etc. |
| **Price** | Free/Paid |

---

*[Back to Category](../)*
```

## Development

### Prerequisites

- Go 1.24+
- Docker (for containerized deployment)

### Running Locally

```bash
# Clone the repository
git clone https://github.com/uberswe/awesomehytale.com.git
cd awesomehytale.com

# Install dependencies
go mod download

# Run the server
go run main.go

# Open http://localhost:8080
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `BASE_URL` | `http://localhost:8080` | Base URL for canonical URLs |
| `ENVIRONMENT` | `development` | Environment (production/development) |

### Building Docker Image

```bash
docker build -t awesomehytale .
docker run -p 8080:8080 awesomehytale
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This website is not affiliated with, endorsed by, or connected to Hypixel Studios or Hytale. "Hytale" is a trademark of Hypixel Studios.

## Contact

- Website: [awesomehytale.com](https://www.awesomehytale.com)
- Email: hello@awesomehytale.com
- GitHub: [github.com/uberswe/awesomehytale.com](https://github.com/uberswe/awesomehytale.com)
