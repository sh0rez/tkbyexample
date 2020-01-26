module.exports = {
  siteMetadata: {
    title: `Tanka by Example`,
    description: `Tanka examples`,
    author: `@sh0rez`,
  },
  plugins: [
    {
      resolve: `gatsby-source-filesystem`,
      options: {
        path: `${__dirname}/dist`,
        name: `examples`,
      },
    },
    `gatsby-plugin-react-helmet`,
    {
      resolve: `gatsby-transformer-remark`,
      options: {
        plugins: [
          {
            resolve: `gatsby-remark-vscode`,
            options: {
              colorTheme: "Light+ (default light)",
              injectStyles: false,
              extensions: [
                {
                  identifier: "heptio.jsonnet",
                  version: "0.1.0",
                },
              ],
            },
          },
        ],
      },
    },
    `gatsby-plugin-catch-links`,
    {
      resolve: `gatsby-plugin-manifest`,
      options: {
        name: "Tanka by Example",
        short: "tkx",
        icon: "src/img/tkx.png",
        start_url: "/",
        display: "standalone",
        background_color: "#ffffff",
        theme_color: "#000000",
      },
    },
  ],
}
