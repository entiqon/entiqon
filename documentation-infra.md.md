# 🧾 Entiqon Documentation Workflow (Docker + GitHub Pages)

This project uses a **fully Dockerized MkDocs setup** with GitHub Actions for clean and automatic documentation deployment.

---

## 🧠 Overview

- ✅ No `/docs-site/` or `gh-pages` branch required
- ✅ MkDocs runs **inside Docker** in GitHub Actions
- ✅ All documentation files are kept in the main branch
- ✅ GitHub Pages is automatically updated on every push to `main`

---

## 📁 Project Structure

```
/mkdocs.yml                # MkDocs config file
/Dockerfile                # Docker image for building docs
/requirements.txt          # For optional local usage or reference
/docs/                     # Additional architecture docs (optional)
/releases/                 # Changelog and release notes
/*.md                      # Builder guides, index, overview
/.github/workflows/docs.yml  # GitHub Actions deploy pipeline
```

---

## 🚀 How It Works

1. **Docker container is built on GitHub Actions** from `Dockerfile`
2. **MkDocs builds** the site from all referenced `.md` files
3. **Output is mounted** to `/site` and uploaded as a GitHub Pages artifact
4. **Deployment happens** via GitHub's native `pages` action

---

## 🛠 Local Preview (Optional)

You don't need Python or pip — just Docker:

```bash
docker build -t entiqon-docs .
docker run -p 8000:8000 entiqon-docs
```

Access: [http://localhost:8000](http://localhost:8000)

---

## 📦 Publishing Workflow

Triggered on every push to `main`.

```yaml
.github/workflows/docs.yml
```

It runs:
- `docker build`
- `mkdocs build`
- Uploads `/site` to GitHub Pages

---

## 🧪 Maintainers

To update docs:
- Edit Markdown files or `mkdocs.yml`
- Commit and push to `main`
- GitHub Actions will take care of the rest

No manual `mkdocs build`, no clutter.