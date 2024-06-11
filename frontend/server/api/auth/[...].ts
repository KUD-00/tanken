import { NuxtAuthHandler } from '#auth'
import GithubProvider from 'next-auth/providers/github'
import crypto from 'crypto';

const runtimeConfig = useRuntimeConfig()

const generateSecret = () => {
  return crypto.randomBytes(64).toString('hex');
};

export default NuxtAuthHandler({
  secret: generateSecret(),
  providers: [
    // @ts-expect-error You need to use .default here for it to work during SSR. May be fixed via Vite at some point
    GithubProvider.default({
      clientId: runtimeConfig.githubClientId,
      clientSecret: runtimeConfig.githubClientSecret,
    }),
  ],
});
