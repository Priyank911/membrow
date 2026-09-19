import adapter from "@sveltejs/adapter-static";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    // Electron loads the packaged renderer from file://, so generated assets
    // must be relative to index.html instead of rooted at the drive letter.
    paths: {
      relative: true
    },
    adapter: adapter({
      fallback: 'index.html'
    }),
  },
};

export default config;
