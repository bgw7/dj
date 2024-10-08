const backend_api = "./api"
function searchVenues(queryParam) {
  return fetch(backend_api + "/venues/search?q=" + queryParam)
    .then(res => {
      if (res.status == 200) {
        return res.json();
      } else {
        throw Error('API response status not 200')
      }
    })
}
function getVenues() {
  return fetch(backend_api + "/venues/")
    .then(res => {
      if (res.status == 200) {
        return res.json();
      } else {
        throw Error('API response status not 200')
      }
    })
}

function createVenue(obj) {
  let body = JSON.stringify(obj)
  return fetch(backend_api + "/reservation/", {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    },
    body: body,
    cache: 'default'
  })
    .then(res => {
      if (res.status == 200) {
        return res.json();
      } else {
        throw Error('API response status not 200')
      }
    })
}

export { getVenues, createVenue, searchVenues }
