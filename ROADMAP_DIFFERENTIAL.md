# Roadmap and Differential for eks-review-cli

This document outlines strategies and planned features that differentiate **eks-review-cli** from standard tools like `kubectl`.

## 🎯 Differential Analysis: How to Stand Out?
While `kubectl` is the default tool to interact with Kubernetes, `eks-review-cli` aims to add value by focusing on:

1. **Insights and Analysis**
   - Go beyond raw listings of resources.
   - Provide smart summaries that highlight problems or areas that need attention.
   - Detect common issues and offer suggestions.
   - Correlate information from different sources for richer context.

2. **Specific focus on EKS**
   - Include checks and listings unique to Amazon EKS (control plane configuration, node groups, IRSA, CNI, etc.).
   - Validate against AWS best practices for EKS.

3. **Optimization and Cost Management**
   - Identify underused or unused resources (ConfigMaps, Secrets, PVCs, PVs).
   - Provide recommendations on sizing `requests` and `limits`.
   - (Advanced) Suggest cost optimizations on nodes.

4. **Proactive Security**
   - Basic RBAC audits for overly permissive roles.
   - Analyze NetworkPolicies for risky configurations.
   - Check security settings in images or ServiceAccounts.

5. **Improved User Experience**
   - Guided troubleshooting flows.
   - Friendlier output with colors or actionable hints.
   - Customizable review profiles.

The goal is not to replace `kubectl` but to complement it, offering intelligence, automating common review tasks and focusing on EKS.

---

## ✅ Development Checklist and Differential Features
Below is a prioritized list of potential features.

### Phase 1: Monitoring Improvements and Quick Insights (Short Term)
- [ ] Enhanced `monitor status` with colored output highlighting problems.
- [ ] Quick alert section summarizing the most critical issues.
- [ ] New command `monitor health-check` (or `monitor quick-insights`) for basic checks:
  - Nodes in `NotReady` state.
  - Pods with problematic states (`Failed`, `CrashLoopBackOff`, `ImagePullBackOff`, long `Pending`).
  - Workloads missing ready replicas.
  - Recent `Warning` events.
- [ ] Improved `monitor get pods` with an option to show last warning events.
- [ ] Documentation updates for new commands.

### Phase 2: Deeper Analysis and EKS Focus (Medium Term)
- [ ] Enriched output for `monitor get <resource>`; for deployments show ReplicaSet and Pod status.
- [ ] Investigate AWS SDK Go v2 APIs to obtain EKS control plane and NodeGroup info.
- [ ] New command `eks-review eks status` (or `monitor eks-info`).
- [ ] New command `eks-review eks check-irsa <namespace>/<serviceaccount>`.
- [ ] Basic optimize module: commands to list unused ConfigMaps and Secrets.

### Phase 3: Security and Advanced Optimization (Long Term)
- [ ] Security module commands such as `security overly-permissive-sa` and `security open-networkpolicies`.
- [ ] Diagnose module with `diagnose pod <pod_name>` summarizing description, logs, events and container status.
- [ ] General improvements like colored output and review profiles.

---

**Next Steps Suggestion**
1. Start with Phase 1 features such as the improved `monitor status` or the new `monitor health-check`.
2. Choose a small, manageable improvement to gain traction—for example, highlighting pods in bad state.
3. Iterate, test yourself and consider enhancements.

Use this roadmap as a guide and adapt it according to your interests and feedback. The aim is to build something useful and differentiated!
