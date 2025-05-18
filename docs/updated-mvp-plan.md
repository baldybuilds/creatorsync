# CreatorSync MVP Plan

## 🧠 What is CreatorSync?

CreatorSync is a fullstack web app designed to help Twitch streamers repurpose their content for platforms like TikTok, YouTube Shorts, and more — all from a single dashboard.

It provides creators with:

- Seamless login via Twitch, Google, or Discord
- A dashboard showing their most recent Twitch streams and clips
- An **in-browser video editor** (powered by React Video Editor + Remotion)
- Server-side video rendering using AWS Lambda
- Simple export and download of final clips

---

## 🎯 Why Are We Building It?

Streamers today face 3 key problems:

1. **Time-wasting workflow** — manually downloading, trimming, editing, and re-uploading content
2. **Inconsistency** — struggling to maintain a regular social media posting schedule
3. **Tool overload** — using 5+ separate apps to manage what should be a simple content loop

### ✅ CreatorSync Solves This By:

- Providing a **centralized dashboard** for stream content
- Enabling quick, visual editing in the browser
- Handling the entire rendering/export process for them
- (Eventually) offering one-click publishing to platforms

---

## 📦 MVP Goals

We want to launch a lean but functional MVP that allows a small number of real streamers to:

- Log in
- View clips
- Edit clips
- Export and download their clip

At least **5 users exporting 3+ clips within 2 weeks** will be our MVP success indicator.

---

## 🔧 Tech Stack Overview

- **Frontend**: React (via Go-Blueprint) + Tailwind CSS V4 + Clerk
- **Backend**: Go + Fiber
- **Database**: PostgreSQL on Railway
- **Video Rendering**: Remotion + AWS Lambda
- **Storage**: AWS S3 (for rendered videos) - _or local storage for MVP_
- **Auth**: Clerk (Google, Discord, Twitch)
- **Analytics & Feature Flags**: PostHog
- **Build Process**: Go-Blueprint scaffolding

---

## 🧭 Milestone Plan

| Milestone           | Target Date | Description                           |
| ------------------- | ----------- | ------------------------------------- |
| ✅ Project scaffold | Week 1      | Go-Blueprint setup, PostgreSQL, Clerk |
| 🔐 Auth + Clips     | Week 2      | Login + fetch Twitch clips            |
| ✂️ Editor Ready     | Week 3      | React Video Editor w/ dummy clips     |
| 🎬 Rendering Setup  | Week 4      | Lambda rendering + file storage       |
| ✅ MVP Soft Launch  | Week 5      | 20 testers invited                    |

---

## 📌 Architecture Notes

- **Frontend**: React SPA built with Go-Blueprint, served by Go backend
- **Rendering Pipeline**: React client → Go backend validates → Remotion Lambda renders → Storage saves → Database tracks status
- **Authentication**: Clerk handles OAuth, Go backend validates Clerk JWT tokens
- **Asset Flow**: Twitch clips → React Video Editor → Trim parameters → Remotion rendering
- **Database**: PostgreSQL for both application data and render job statuses

---

## 🚀 Key Implementation Decisions

1. **Why Go-Blueprint?** Provides structured fullstack setup with React + Go, handles build process automatically
2. **Why Go + Fiber?** Fast API development, excellent for handling multiple render jobs concurrently
3. **Why Clerk for Auth?** Pre-built Twitch OAuth integration with proper scopes management
4. **Why Railway PostgreSQL?** Simple deployment, good for MVP, no need for separate backend services
5. **Why Remotion Lambda?** Offloads heavy video processing without blocking the main application

---

## 🗂️ Storage Strategy Options

Since AWS S3 is new to you, consider these MVP approaches:

### Option 1: Start Simple (Recommended)

- Store rendered videos temporarily on server filesystem
- Return download links that expire after 24 hours
- Migrate to S3 post-MVP for scalability

### Option 2: S3 with Minimal Setup

- Use AWS SDK for Go (`aws-sdk-go-v2`)
- Single S3 bucket with public read permissions
- Basic upload/download functionality

### Option 3: Railway Volume Storage

- Use Railway's persistent volumes
- Similar to local filesystem but with better persistence
- Good middle ground before S3 migration

---

## 📈 Success Metrics

- User onboarding completion rate > 80%
- Average time from clip selection to download < 5 minutes
- Render success rate > 95%
- User retention: 60% of users create 2nd clip within week 1

---

## 📝 Storage Implementation Notes

For MVP simplicity, I recommend starting with **local/temporary storage**:

```go
// Pseudocode for MVP storage approach
func handleRenderedVideo(jobID string, videoBytes []byte) error {
    // Store temporarily on server
    filepath := fmt.Sprintf("/tmp/renders/%s.mp4", jobID)
    ioutil.WriteFile(filepath, videoBytes, 0644)

    // Update database with download link
    downloadURL := fmt.Sprintf("/api/download/%s", jobID)
    db.UpdateJob(jobID, "completed", downloadURL)

    // Clean up file after 24 hours
    time.AfterFunc(24*time.Hour, func() {
        os.Remove(filepath)
    })
}
```

Post-MVP, you can migrate to S3 without changing the frontend interface.
