import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { DataFetcherService } from "~/rpc/data-fetcher-service_connect";
import { getServerSession } from "#auth";

const transport = createConnectTransport({
  baseUrl: "http://data-fetcher:50051",
  httpVersion: "2",
});

const client = createPromiseClient(DataFetcherService, transport);

export default defineEventHandler(async (event) => {
  const session = await getServerSession(event);
  if (!session) {
    return { status: "unauthenticated!" };
  }

  const { method } = await readBody(event);
  if (method === "uploadNewPost") {
    const { userId, pictureChunk, location, content, tags } = await readBody(
      event
    );

    if (session.user?.name !== userId) {
      return { status: "unauthorized!", message: "User ID mismatch!" };
    }

    try {
      const response = await client.addPost({
        userId,
        pictureChunk: new Uint8Array(pictureChunk),
        location,
        content,
        tags,
      });

      console.log("Post uploaded successfully:", response);
      return response;
    } catch (error) {
      console.error("Failed to upload post:", error);
    }
  }
});
