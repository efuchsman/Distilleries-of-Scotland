
<h1 align="center">Distilleries of Scotland API</h1>

<br />
<div align="center">
  <a href="https://github.com/efuchsman/Distilleries-of-Scotland">
    <img src="https://www.edinburghwhiskyacademy.com/cdn/shop/articles/EWA_Scotch-whisky-regions_1200x1200.webp?v=1689791862" alt="Map of the Distillery Regions in Scotland" width="650" height="450">
  </a>

  <h3 align="center">
    Get regions and distilleries within each region!
    <br />
  </h3>
</div>

# Table of Contents
* [App Description](#app-description)
* [Learning Goals](#learning-goals)
* [System Requirements](#system-requirements)
* [Technologies Used](#technologies-used)
* [Setup](#setup)
* [Respository](#repository)
* [Endpoints](#endpoints)
* [Acknowledgements](#acknowledgments)
* [Contact](#contact)

# App Description:
* Distilleries of Scotland allows users to search for the Scottish regions, and the distilleries within each of those regions.

# Learning Goals:
* Build Rest API in Go from scratch exposing endpoints for regions, region, and regional distilleries
* Host database in PostgreSQL
* Build Environment Configurations
* Containerize
* Test DB calls using single transactional databases
* Build interfaces for mocking
* Build and deploy locally with kubernetes config files using MiniKube
* Allow CORS for a potential Frontend integration
* CI/CD with CircleCI

# System Requirements
* go 1.21
* Docker Desktop
* Minikube
* Postgres will be installed using `startup.sh`

# Technologies used
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
![Docker](https://img.shields.io/badge/Docker-ffffff?style=for-the-badge&logo=docker&logoColor=ffffff&color=0db7ed)
[![Minikube](https://img.shields.io/badge/Minikube-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)](https://minikube.sigs.k8s.io/)
[![CircleCI](https://img.shields.io/badge/CircleCI-343434?style=for-the-badge&logo=circleci&logoColor=white)](https://circleci.com/)

# Setup
- run: `go mod tidy`
- run: `chmod +x startup.sh`
- From the root of the directory, run: `./startup.sh`
- To use MiniKube, first run `miniKube start --profile=distilleries-of-scotland`, then run: `chmod +x deploy.sh`, followed by `./deploy.sh`
- To see the API on localhost, run: `kubectl port-forward service/distilleries-of-scotland-service 8000:8000`

- To run locally without MiniKube, run: `go run distilleries_of_scotland.go`

# Repository

https://github.com/efuchsman/Distilleries-of-Scotland

# Endpoints

## GET All Regions:

`GET '/regions'`

```
  {
      "regions": [
          {
              "region_name": "Speyside",
              "description": "The most densely populated Whisky region in the world, famous for fertile glens and, of course, the River Spey. Speyside whiskies are known for being frugal with peat and full of fruit. Apple, pear, honey, vanilla and spice all have a part a role in expressions from this region, which are commonly matured in Sherry casks."
          },
          {
              "region_name": "Lowland",
              "description": "Soft and smooth malts are characteristic of this region, offering a gentle, elegant palate reminiscent of grass, honeysuckle, cream, ginger, toffee, toast and cinnamon. The whiskies are often lighter in character and perfect for pre-dinner drinks."
          },
          {
              "region_name": "Highland",
              "description": "This region, which also takes in the islands, has a huge diversity of flavours and characters. From lighter whiskies all the way through salty coastal malts, the Highlands offers a Scotch for all palates."
          },
          {
              "region_name": "Campbeltown",
              "description": "Campbeltown whiskies are varied and full of flavour. Hints of salt, smoke, fruit, vanilla and toffee mingle in whiskies of robust and rich character."
          },
          {
              "region_name": "Islay",
              "description": "Islay (pronounced ‘eye-luh’) is a magical island where the majority of its population are involved in whisky production. Famous for fiery, heavily peated whiskies."
          }
      ]
  }
```

## GET a Region

`GET 'regions/:region_name'`

```
{
    "region_name": "Speyside",
    "description": "The most densely populated Whisky region in the world, famous for fertile glens and, of course, the River Spey. Speyside whiskies are known for being frugal with peat and full of fruit. Apple, pear, honey, vanilla and spice all have a part a role in expressions from this region, which are commonly matured in Sherry casks."
}
```
## GET Disitlleries for a Region

`GET 'regions/:region_name/distilleries`

```
{
    "distilleries": [
        {
            "distillery_name": "Glen Scotia",
            "region_name": "Campbeltown",
            "geo": "55.429514, -5.604236",
            "town": "Campbeltown",
            "parent_company": "Loch Lomond Group"
        },
        {
            "distillery_name": "Glengyle",
            "region_name": "Campbeltown",
            "geo": "55.427208, -5.61095",
            "town": "Campbeltown",
            "parent_company": "J. & A. Mitchell and Co."
        },
        {
            "distillery_name": "Springbank",
            "region_name": "Campbeltown",
            "geo": "55.425417, -5.608861",
            "town": "Campbeltown",
            "parent_company": "J. & A. Mitchell and Co."
        }
    ]
}
```

# Contact

<table align="center">
  <tr>
    <td><img src="https://avatars.githubusercontent.com/u/104859844?s=150&v=4"></td>
  </tr>
  <tr>
    <td>Eli Fuchsman</td>
  </tr>
  <tr>
    <td>
      <a href="https://github.com/efuchsman">GitHub</a><br>
      <a href="https://www.linkedin.com/in/elifuchsman/">LinkedIn</a>
   </td>
  </tr>
</table>



# Acknowledgements
  - https://www.scotch-whisky.org.uk/discover/enjoying-scotch/scotch-whisky-regions/

  - https://www.edinburghwhiskyacademy.com/cdn/shop/articles/EWA_Scotch-whisky-regions_1200x1200.webp?v=1689791862


