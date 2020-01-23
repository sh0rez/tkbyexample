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
  ],
}
