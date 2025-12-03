# 🧩 Custom Kubernetes Controller (Learning Project)

This repository contains a basic custom Kubernetes controller written in Go, built using `client-go` informers and workqueues. This is a learning-focused project created to understand how Kubernetes controllers work internally.

The controller watches Deployment resources and automatically creates a Service for each Deployment.

## 🚀 What This Project Does

* Connects to a Kubernetes cluster using:
    * `kubeconfig` (local)
    * `InClusterConfig` (for in-cluster usage)
* Uses Shared Informer Factory
* Watches Deployment objects
* On Deployment creation, it:
    * Extracts pod labels
    * Automatically creates a Service
* Uses:
    * Workqueue
    * Lister
    * Cache sync
* Implements a basic controller worker loop

**In short:**

📦 Whenever a Deployment is created, a matching Service is auto-created.

## 🛠 Tech Stack

* Go (Golang)
* Kubernetes
* client-go
* Informers
* Workqueue
* Listers

## 📁 Project Structure

```
.
├── main.go
├──chart
│  └── deployment.yaml
├── controller/
│   └── controller.go
└── README.md
```

| File | Purpose |
|------|---------|
| `main.go` | Initializes Kubernetes client, informers & controller |
| `controller/controller.go` | Core controller logic |
| `README.md` | Project documentation |

## ⚙️ How It Works (Simple Flow)

1. Kubernetes cluster connection is created
2. Deployment informer is initialized
3. Event handlers are attached:
    * `Add`
    * `Delete`
4. When a Deployment is added:
    * It is pushed to the workqueue
5. Worker picks it up and:
    * Reads Deployment from cache
    * Creates a matching Service

## ▶️ How to Run This Project

### ✅ Prerequisites

* Go installed
* Kubernetes cluster (local or remote)
* `kubectl` configured

### ▶️ Run Locally Using Kubeconfig

```bash
go run main.go --kubeconfig=/path/to/.kube/config
```

### ▶️ Create a Test Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: ekposetest
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
```

Once applied:

✅ Your controller will automatically create a matching Service.

## 📚 What I Learned From This Project

* How Kubernetes controllers work internally
* Difference between:
    * Informers
    * Listers
    * Workqueues
* How caching works with `WaitForCacheSync`
* How event-driven systems work in Kubernetes
* How to:
    * Watch resources
    * Process them asynchronously
    * Perform reconciliation
* How Service selectors work from Deployment labels

## ⚠️ Current Limitations (Intentionally Left for Learning)

* No Update handler yet
* No retry logic with rate limiting
* No ownership reference between Deployment & Service
* Service already exists error is not handled fully
* No leader election

These are intentionally left incomplete as next learning milestones.

## 🎯 Why This Repository Exists

This is not a production project. It exists to:

* Track my Kubernetes controller learning journey
* Show real hands-on implementation

## ✅ Next Improvements (Planned)

* Add `UpdateFunc`
* Add retry handling with `AddRateLimited`
* Add OwnerReferences
* Handle Service already exists properly
* Convert into a full Operator using Kubebuilder

## 🙌 Final Note

This repository reflects:

* My hands-on learning with Kubernetes controllers
* My exploration of platform engineering internals
* My consistency in pushing meaningful code to GitHub

If you're reviewing this as a recruiter or mentor .... feedback is always welcome 😊