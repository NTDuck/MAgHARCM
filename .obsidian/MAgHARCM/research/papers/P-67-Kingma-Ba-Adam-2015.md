---
title: "P-67 — Kingma & Ba 2015 — Adam: A Method for Stochastic Optimization"
backlink: "[[1.0.0 P-67]]"
aliases:
  - "1.0.0 P-67"
  - "P-67"
  - "P-67-Kingma-Ba-Adam-2015"
  - "P-67-Kingma-Ba-Adam-2015"
  - "Kingma-Ba-Adam-2015"
tags: [paper, optimizer, adam, deep-learning, transfer-learning, [[1.0.0 P-66]], hop-2]
---

# [[1.0.0 P-67 — Adam Optimizer]]

## Citation

Kingma, D. P., & Ba, J. (2015). *Adam: A Method for Stochastic Optimization*. ICLR 2015. arXiv:1412.6980. URL: https://arxiv.org/abs/1412.6980.

## Summary

Adam combines the **momentum** method (smoothed first-moment of gradients) with **RMSProp** (adaptive per-parameter learning rates from second-moment estimates). The update rule is:

$$m_t = \beta_1 m_{t-1} + (1 - \beta_1) g_t$$
$$v_t = \beta_2 v_{t-1} + (1 - \beta_2) g_t^2$$
$$\hat{m}_t = m_t / (1 - \beta_1^t), \quad \hat{v}_t = v_t / (1 - \beta_2^t)$$
$$\theta_t = \theta_{t-1} - \alpha \hat{m}_t / (\sqrt{\hat{v}_t} + \epsilon)$$

with defaults $\beta_1 = 0.9$, $\beta_2 = 0.999$, $\epsilon = 10^{-8}$. The bias-correction terms $\hat{m}_t, \hat{v}_t$ are critical for the early steps when $m_0 = v_0 = 0$.

Adam is the **default optimizer** for nearly every transformer training run since 2017 (BERT, GPT-2, RoBERTa, T5, LLaMA, Qwen). It converges faster than SGD on noisy gradients and rarely requires learning-rate warmup beyond a brief initial period.

## Method

Adam's contribution is the **bias correction** that makes the first-moment and second-moment estimates unbiased when initialized at zero. Without it, the early training steps have systematically underestimated moments, leading to oversized updates that destabilize training. The correction term $1 - \beta^t$ decays to zero as $t$ grows, so the correction is purely an initialization fix.

## Findings Relevant to MAgHARCM

- **Default optimizer for every LoRA / domain-adapter training** on Qwen2.5-Coder. MAgHARCM's per-domain adapter recipe (cf. [[1.0.0 P-58]], [[1.0.0 P-66]]) uses AdamW (Adam with decoupled weight decay) at $\beta_1 = 0.9$, $\beta_2 = 0.999$, $\epsilon = 10^{-8}$, learning rate $2 \times 10^{-4}$ with linear warmup over 100 steps.
- **Adam's bias-correction** is the empirical justification for the **warmup-then-decay** learning-rate schedule: the early-step bias correction gives the warmup a self-stabilizing property; Adam with no warmup can diverge on transformer gradients.
- **Adam vs. SGD for SLM fine-tuning**: Adam converges in ~10× fewer steps than SGD on small datasets, which is the critical property for MAgHARCM's per-fragment training samples.

## How MAgHARCM Uses It

`compiletime.DefaultAdamBeta1`, `compiletime.DefaultAdamBeta2`, `compiletime.DefaultAdamEpsilon`, `compiletime.DefaultAdamLearningRate` constants in `internal/compiletime/compiletime.go` define the Adam hyperparameters for any domain-adapter training. The constants live in the centralised compile-time config (per directive 4: hardcoded magic values go to a centralised config).

## References

### Hop-1 (Kingma & Ba 2015 cites)
- Duchi, J., Hazan, E., Singer, Y. (2011). *Adaptive Subgradient Methods for Online Learning and Stochastic Optimization*. JMLR 12:2121-2159. (AdaGrad)
- Tieleman, T. & Hinton, G. (2012). *Lecture 6e: RMSProp*. Coursera Neural Networks for ML. (RMSProp)
- Polyak, B. T. (1964). *Some methods of speeding up the convergence of iteration methods*. USSR Computational Mathematics and Mathematical Physics 4(5):1-17. (momentum)
- Nesterov, Y. (1983). *A method for unconstrained convex minimization problem with the rate of convergence $O(1/k^2)$*. Soviet Mathematics Doklady 27:372-376.

### Hop-2
- Loshchilov, I. & Hutter, F. (2019). *Decoupled Weight Decay Regularization* (AdamW). ICLR 2019. arXiv:1711.05101.
- Howard, J. & Ruder, S. (2018). See [[1.0.0 P-66]].
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762.

## Backlinks

[[1.0.0 P-58]], [[1.0.0 P-66]], [[1.0.0 P-55]], [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]].

P-67 is the **default optimizer** for MAgHARCM's per-domain LoRA adapter training on Qwen2.5-Coder.
