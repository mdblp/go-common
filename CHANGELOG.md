# Go-Common

Go-Common is the library of common functions for Tidepool's Go-based applications

## 1.6.0 - 2022-09-27
### Engineering
- Remove seagull client (moved to seagull repo)

## 1.5.0 - 2022-08-23
### Engineering
- Add OTP package to this common library

## 1.4.1 - 2022-06-29
### Engineering
- YLP-1644: possibility to set a custom CA cert for when we run backloops in our dev/test env

## 1.4.0 - 2022-06-03
### Added
- Add setCollection method to seagull client

## 1.3.3 - 2022-06-03
### Fixed
- Fix casting issue on auth0 custom claims

## 1.3.2 - 2022-05-31
### Engineering
- Change AUTH0_DOMAIN by AUTH_URL to allow clients to select http vs https

## 1.3.1 - 2022-05-18
### Fixed
- Correct Auth client and its mock

## 1.3.0 - 2022-05-12
### Added
- YLP-1533 Create a new client to authenticate requests based on tidepool tokens or new bearer token

## 1.2.0 - 2021-12-08
### Changed
- YLP-1111 OPA client should handle additional data in requests

## 1.1.0 - 2021-08-19
### Changed
- YLP-956: migrate log to lorgus libraries

## 1.0.0 - 2021-08-11
### Engineering
- YLP-923 Remove hakken & highwater code

## 0.9.0 - 2021-06-09
### Added
- Adding read preferences as config on mongo client

## 0.8.1 - 2021-03-23
### Engineering
- Fixing publish step in jenkins pipeline

## 0.8.0 - 2021-03-16
### Changed
- Removing gatekeeper client
- Removing hakken hostgetter for seagull client
- Enrich seagull mock
### Engineering
- Jenkins pipeline

## 0.7.2 - 2021-03-15
### Fixed
- OPA client: decode query string before sending to OPA

### Changed
- YLP-505: add hcp role, same as done for shoreline change

## 0.7.0 - 2021-03-05
### Added
- YLP-469 Implement Authorization Client for go services

## 0.6.2 - 2020-10-29
### Fixed
- YLP-255 MongoDb context cancellation issue

## 0.6.1 - 2020-09-23
### Fixed
- Fixing mdblp vs tidepool-org import path 

## 0.6.0 - 2020-09-23
### Changed
- PT-1479 Update mongoDb client to be able to start without the database
- PT-1514 Update shoreline client to be able to start without acquiring server token

## 0.5.0 - 2020-06-19
### Added
- PT-1383 Add portal-api client to fetch the PatientConfig

## 0.4.0 - 2020-04-14
### Added
- Complete Getekeeper client with missing route

## 0.3.0 - 2019-10-17
### Changed
- PT-727 Add the versioning info in the status object

## 0.2.0 - 2019-10-09
### Added
- PT-582 Merge from upstream Tidepool v0.4.1

## 0.1.0 - 2019-05-29
### Added
- Merge from upstream Tidepool
