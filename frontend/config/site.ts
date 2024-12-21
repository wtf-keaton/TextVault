export type SiteConfig = typeof siteConfig;

export const siteConfig = {
  name: "TextVault",
  description: "Make beautiful websites regardless of your design experience.",
  navItems: [
    {
      label: "New paste",
      href: "/",
    },
  ],
  navMenuItems: [
    {
      label: "Profile",
      href: "/profile",
    },
    {
      label: "Logout",
      href: "/logout",
    },
  ],
  links: {
    github: "https://github.com/wtf-keaton/",
    githubProject: "https://github.com/wtf-keaton/TextVault",
    login: "/login",
  },
};
