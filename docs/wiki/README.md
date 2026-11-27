# Project Wiki content

This folder holds the pages of the **GitHub Wiki**, versioned alongside the
code. On GitHub the Wiki is a separate git repository (`<repo>.wiki.git`); each
`.md` file here becomes a page and the `images/` folder (optional) holds images.

The Wiki is a **usage guide for the tool**, with one page per command / use
case, in **English** and **Portuguese**.

## Structure

| File | Page |
|------|------|
| `Home.md` | Landing page (language selector) |
| `_Sidebar.md` | Side navigation (EN + PT) |
| `Installation-EN.md` / `Installation-PT.md` | Install / build |
| `Global-Options-EN.md` / `Global-Options-PT.md` | Global flags, DNS server, proxy, timeouts, scoping |
| `Recon-EN.md` / `Recon-PT.md` | `recon` — DNS record enumeration |
| `Brute-Force-EN.md` / `Brute-Force-PT.md` | `brute` — wordlist brute-force |
| `Resolve-EN.md` / `Resolve-PT.md` | `resolve file` / `bloodhound` / `crtsh` |
| `Threat-Analysis-EN.md` / `Threat-Analysis-PT.md` | `threat-analysis` — typosquatting & malicious domains |
| `AD-DNS-Dump-EN.md` / `AD-DNS-Dump-PT.md` | `ad` — ADIDNS dump over LDAP |
| `Wordlist-EN.md` / `Wordlist-PT.md` | `wordlist` — build a clean label wordlist |
| `Reports-EN.md` / `Reports-PT.md` | `report convert` / `report elastic` |
| `Output-and-Writers-EN.md` / `Output-and-Writers-PT.md` | Output formats, writers, SOA linking, control DB |

> GitHub Wiki naming convention: hyphens in the file name become spaces in the
> page title, and links use the name **without** the `.md` extension
> (e.g. `[Recon](Recon-EN)`). Images are referenced as `images/file.png`.

## Automatic publishing (GitHub Actions)

The workflow [`.github/workflows/publish-wiki.yml`](../../.github/workflows/publish-wiki.yml)
publishes this folder to the Wiki **automatically** on every push to `main` that
changes `docs/wiki/**`. It mirrors the content (except this `README.md`) to
`<owner>/<repo>.wiki.git`, removing pages that no longer exist here.

**One-time prerequisite:** enable the Wiki in *Settings → Features → Wikis* and
create **one first page** through the web UI (that initializes the `.wiki.git`
repository). You can also trigger it manually in **Actions → Publish Wiki → Run
workflow**.

## Manual publishing (alternative)

```bash
gh repo clone helviojunior/enumdns.wiki /tmp/enumdns-wiki
rsync -a --delete --exclude='.git/' --exclude='README.md' docs/wiki/ /tmp/enumdns-wiki/
cd /tmp/enumdns-wiki && git add -A && git commit -m "Update Wiki" && git push
```
