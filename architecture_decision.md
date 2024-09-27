# Notification Service Architecture Decision

## Overview

We have a main service and a notification service. When a user's post is liked by another user, we need to send a notification. There are two possible approaches to achieve this:

1. **API Call to the Notification Service**
2. **Push Event to Kafka, and let the Notification Service consume it**

## Option 1: API Call to Notification Service

### Pros:
- **Simplicity**: Direct and straightforward. The main service triggers the notification service by calling an API whenever an event (e.g., post liked) occurs.
- **Real-time**: Notifications are delivered immediately as the API call is made.
- **Less Overhead**: No need to set up and manage additional infrastructure like Kafka.

### Cons:
- **Tight Coupling**: The main service is tightly coupled with the notification service. If the notification service is down or slow, it might impact the main service.
- **Scalability**: With an increasing number of likes, the notification service must handle many requests in real-time, potentially becoming a bottleneck.
- **Error Handling**: Retries and failure management need to be manually handled in case the API call fails.

## Option 2: Event-Driven with Kafka

### Pros:
- **Decoupling**: The main service does not directly interact with the notification service. It simply pushes an event to Kafka, and the notification service processes it asynchronously.
- **Scalability**: Kafka can handle a large volume of events and scale efficiently. The notification service can consume events at its own pace, reducing bottlenecks.
- **Reliability**: Kafka persists events, allowing for retries and ensuring no notifications are lost, even if the notification service is temporarily down.
- **Asynchronous Processing**: The main service focuses on its core operations without being blocked by notification sending.

### Cons:
- **Infrastructure Overhead**: Kafka introduces additional infrastructure that needs to be set up and managed.
- **Eventual Consistency**: Notifications may be delayed due to the asynchronous nature of Kafka-based processing.
- **Development Complexity**: Implementing Kafka producers and consumers adds complexity, including managing topics, partitions, and offsets.

## Decision Points

- **Real-time Requirements**: If immediate notifications with low latency are essential, the **API call** is a better choice. If slight delays are acceptable, Kafka offers better scalability.
- **Scalability**: For high traffic scenarios, Kafka provides better scalability and reliability.
- **Service Decoupling**: If services should operate independently and asynchronously, Kafka is preferable.
- **Failure Handling**: Kafka provides built-in persistence and retry mechanisms, while API calls need manual error handling and retries.

## Recommendation

- **If scalability and decoupling are important**, go with **Kafka**. It's more resilient and scalable, allowing for better long-term maintainability as the system grows.
- **If simplicity is key and the system is small**, the **API call** approach may be more efficient, especially for real-time notification needs.

Since Kafka is already used in other parts of the system, extending its usage for notifications might be easier.
