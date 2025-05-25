# 🧾 Documentation Infrastructure (MkDocs + Docker + GitHub Actions)

This document describes how Entiqon's documentation is built and deployed using a fully containerized MkDocs setup with GitHub Pages.

---

## 🧠 Overview

- ✅ No `/docs-site/` or `gh-pages` branch is required
- ✅ MkDocs runs **inside Docker** during GitHub Actions
- ✅ All documentation files are managed from the main branch
- ✅ GitHub Pages is automatically updated on every push to `main`

---

## 📁 Project Structure

```
/mkdocs.yml                # MkDocs config file
/Dockerfile-documentation # Docker image for building docs
/requirements.txt          # Optional for pip-based builds
/docs/                     # Architecture or additional docs (optional)
/releases/                 # Changelog and release notes
/*.md                      # Builder guides, index, overview
/.github/workflows/docs.yml  # GitHub Actions deploy pipeline
/documentation-infra.md    # ← This file
```

---

## 🚀 Workflow Summary

1. **Docker container is built** on GitHub Actions using `Dockerfile-documentation`
2. **MkDocs builds** the site from referenced `.md` files
3. **Output is mounted** to `/site` and uploaded to GitHub Pages

No temporary folders. No pollution. No manual builds.

---

## 🧪 Local Preview (Optional)

```bash
docker build -t entiqon-docs -f Dockerfile-documentation .
docker run -p 8000:8000 entiqon-docs
```

Visit: [http://localhost:8000](http://localhost:8000)

---

## 📦 GitHub Deployment

> 📄 For full release history, see the [CHANGELOG](CHANGELOG.md).

GitHub Actions automatically deploys the site on push to `main`.

Workflow:
```yaml
.github/workflows/docs.yml
```

Steps:
- `docker build`
- `mkdocs build`
- Upload `/site`
- GitHub Pages deployment

---

## 🌐 GitHub Pages Configuration

To ensure GitHub correctly publishes documentation built by the `docs.yml` workflow:

### ✅ Set the publishing source to **GitHub Actions**:

1. Go to **Settings > Pages**
2. In the **Build and deployment** section:
   - Change **Source** from `Deploy from a branch` to `GitHub Actions`
   - Click **Save**

This ensures GitHub publishes the site generated from Docker (inside GitHub Actions) rather than expecting a static `/docs` folder on the main branch.

> 💡 If you are using a custom domain, add a `.github/CNAME` file containing your domain name. GitHub will use this file to apply the custom domain without manual entry.

---

## 🧩 Contributors

To update documentation:
- Edit any Markdown file or `mkdocs.yml`
- Commit to `main`
- GitHub will handle the deployment

This file (`documentation-infra.md`) is part of the deployed docs and should remain up to date with the current infrastructure.