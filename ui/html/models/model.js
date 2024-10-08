export class Reservation {
    id = ""
    clientId = ""
    venueId = ""
    startTimestamp = ""
    endTimestamp = ""
    createdBy = ""
    createdAt = ""
    updatedBy = ""
    updatedAt = ""
}

export class Venue {
  type = VenueType
  zipCode = 0
  city = ""
  state = ""
  image = ""
  imageGallery = []
}

export class VenueType {
    name = ""
}
