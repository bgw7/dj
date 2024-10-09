const backend_api = "./api"
function getTracks() {
  return fetch(backend_api + "/tracks/")
    .then(res => {
      if (res.status == 200) {
        return res.json();
      } else {
        throw Error('getTracks API response status not 200')
      }
    })
}

function createTrack(obj) {
  let body = JSON.stringify(obj)
  return fetch(backend_api + `/tracks/`, {
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
        throw Error('createTrack API response status not 200')
      }
    })
}

function createVote(trackId) {
  return fetch(backend_api + `/tracks/${trackId}/`, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    },
    cache: 'default'
  })
    .then(res => {
      if (res.status == 200) {
        return res.json();
      } else {
        throw Error('createVote API response status not 200')
      }
    })
}

function deleteVote(trackId) {
  return fetch(backend_api + `/tracks/${trackId}/`, {
    method: 'DELETE',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    },
    cache: 'default'
  })
    .then(res => {
      if (res.status == 200) {
        return res.json();
      } else {
        throw Error('deleteVote API response status not 200')
      }
    })
}

export { getTracks, createTrack, createVote, deleteVote }
