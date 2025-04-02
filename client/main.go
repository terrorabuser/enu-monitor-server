package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "golang_gpt/internal/proto"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewContentServiceClient(conn)

	// Timeout for context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Example 1: Add new content
	fmt.Println("Adding new content...")
	addReq := &pb.AddContentRequest{
		UserId:    2,
		Building:  "Building B",
		Floor:     3,
		Notes:     "Conference room",
		FileName:  "test.jpg",
		FilePath:  "/uploads/test.jpg",
		StartTime: timestamppb.Now(),
		EndTime:   timestamppb.New(time.Now().Add(time.Hour * 24)),
	}
	addRes, err := client.AddContent(ctx, addReq)
	if err != nil {
		log.Fatalf("Error during AddContent: %v", err)
	}
	fmt.Printf("Content added with ID: %d\n", addRes.GetContentId())
	contentID := addRes.GetContentId()

	// Example 2: Send content to moderation
	fmt.Println("Sending content to moderation...")
	sendReq := &pb.SendContentToModerationRequest{
		UserId:    2,
		ContentId: contentID,
	}
	sendRes, err := client.SendContentToModeration(ctx, sendReq)
	if err != nil {
		log.Fatalf("Error during SendContentToModeration: %v", err)
	}
	fmt.Printf("Content moderation status: %v, message: %s\n", sendRes.GetStatus(), sendRes.GetMessage())

	// Example 3: Moderate content
	fmt.Println("Moderating content...")
	moderateReq := &pb.ModerateContentRequest{
		UserId:    111, // moderator ID
		ContentId: contentID,
		StatusId:  3,
		Reason:    "Looks good",
	}
	moderateRes, err := client.ModerateContent(ctx, moderateReq)
	if err != nil {
		log.Fatalf("Error during ModerateContent: %v", err)
	}
	fmt.Printf("Moderation success: %v\n", moderateRes.GetSuccess())

	// Example 4: Get contents
	fmt.Println("Getting contents...")
	getReq := &pb.GetContentsRequest{
	}
	getRes, err := client.GetContents(ctx, getReq)
	if err != nil {
		log.Fatalf("Error during GetContents: %v", err)
	}
	fmt.Printf("Found %d contents:\n", len(getRes.GetContents()))
	for _, content := range getRes.GetContents() {
		fmt.Printf("ID: %d, User: %d, File: %s, Status: %d\n",
			content.GetId(),
			content.GetUserId(),
			content.GetFileName(),
			content.GetLatestHistory().GetStatusId())
	}
}