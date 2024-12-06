db.createCollection('quotes')
db.quotes.createIndex({ id: 1 }, { unique: true })
db.quotes.createIndex({ creationDate: 1})
db.quotes.insertMany([
  {
    id: "_dfC0aL_AGD4",
    content: "Great ideas often receive violent opposition from mediocre minds.",
    author: "Albert Einstein",
    authorSlug: "albert-einstein",
    length: 65,
    tags: [
      "Famous Quotes",
      "Technology"
    ],
    creationDate: "2024-03-24"
  },
  {
    id: "U6Al9aA7g7",
    content: "One machine can do the work of fifty ordinary men. No machine can do the work of one extraordinary man.",
    author: "Elbert Hubbard",
    authorSlug: "elbert-hubbard",
    length: 103,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-27"
  },
  {
    id: "DNjQty5jeU",
    content: "Communications tools don't get socially interesting until they get technologically boring.",
    author: "Clay Shirky",
    authorSlug: "clay-shirky",
    length: 90,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "ZyIkYFat1B",
    content: "Computers are like bikinis. They save people a lot of guesswork.",
    author: "Sam Ewing",
    authorSlug: "sam-ewing",
    length: 64,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "APsci40ULi",
    content: "Technology frightens me to death. It's designed by engineers to impress other engineers. And they always come with instruction booklets that are written by engineers for other engineers — which is why almost no technology ever works.",
    author: "John Cleese",
    authorSlug: "john-cleese",
    length: 233,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "Llo063kGTo",
    content: "Ethics change with technology.",
    author: "Larry Niven",
    authorSlug: "larry-niven",
    length: 30,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "KMTJ0Ya3e9",
    content: "To invent, you need a good imagination and a pile of junk.",
    author: "Thomas Edison",
    authorSlug: "thomas-edison",
    length: 58,
    tags: [
      "Imagination",
      "Creativity",
      "Science",
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "yUVQvOdsif",
    content: "Technology made large populations possible; large populations now make technology indispensable.",
    author: "Joseph Wood Krutch",
    authorSlug: "joseph-wood-krutch",
    length: 96,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "qVYnD_eLg5",
    content: "TV and the Internet are good because they keep stupid people from spending too much time out in public.",
    author: "Douglas Coupland",
    authorSlug: "douglas-coupland",
    length: 103,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "JxE5YMTDIK",
    content: "Technology is a word that describes something that doesn't work yet.",
    author: "Douglas Adams",
    authorSlug: "douglas-adams",
    length: 68,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "mQqRWF49Ug",
    content: "The ultimate promise of technology is to make us master of a world that we command by the push of a button.",
    author: "Volker Grassmuck",
    authorSlug: "volker-grassmuck",
    length: 107,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "NoaRFCJNzT",
    content: "So much technology, so little talent.",
    author: "Vernor Vinge",
    authorSlug: "vernor-vinge",
    length: 37,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "_0CfFQ4la0aN",
    content: "If you can't explain it simply, you don't understand it well enough.",
    author: "Albert Einstein",
    authorSlug: "albert-einstein",
    length: 68,
    tags: [
      "Famous Quotes",
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "QQDu56d67y",
    content: "Humanity is acquiring all the right technology for all the wrong reasons.",
    author: "Buckminster Fuller",
    authorSlug: "buckminster-fuller",
    length: 73,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "_hJS3LX4Qz",
    content: "Technology has to be invented or adopted.",
    author: "Jared Diamond",
    authorSlug: "jared-diamond",
    length: 41,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "WSpdlKZYCP",
    content: "Technology… the knack of so arranging the world that we don't have to experience it.",
    author: "Max Frisch",
    authorSlug: "max-frisch",
    length: 84,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  },
  {
    id: "VgsphQQC9g",
    content: "It's supposed to be automatic, but actually you have to push this button.",
    author: "John Brunner",
    authorSlug: "john-brunner",
    length: 73,
    tags: [
      "Technology"
    ],
    creationDate: "2024-03-28"
  }
])
