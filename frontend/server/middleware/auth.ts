import { getToken } from "#auth";

export default defineEventHandler(async (event) => {
  if (event.path.startsWith("/api/rpc")) {
    const token = await getToken({ event });

    if (!token) {
      throw createError({ statusCode: 401, message: "Unauthorized" });
    }

    event.context.user = {
      userId: token.sub,
      accessToken: token.accessToken,
    };
  }
});
