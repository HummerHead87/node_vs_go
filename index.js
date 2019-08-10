const max = require('lodash/max')
const forEach = require('lodash/forEach')
const sortBy = require('lodash/sortBy')

const Binance = require('binance-api-node').default

const client = Binance()

const coinsType = [
  'BTC',
  'BNB',
  'ETH',
  'PAX',
  'USDC',
  'USDT',
  'TUSD',
  'USDS',
  'XRP'
]
getPairs()
// getExchanges()

async function getPairs() {
  const prices = await client.prices()
  
  const date = new Date()
  let pairs = []
  forEach(prices, (priceStr, symbol) => {
    coinsType.some(coin => {
      const regexp = new RegExp(`(.*)(${coin}$)`)
      const matched = symbol.match(regexp)

      if (matched) {
        const pair = {
          pair: [matched[1], matched[2]].join('_'),
          price: parseFloat(priceStr)
        }

        pairs.push(pair)
        return true
      }

      return false
    })
  })

  pairs = sortBy(pairs, o => o.pair)

  console.log('Time: ', new Date - date)
  // console.log(pairs)
  // const keysLength = Object.keys(prices).map(key => key.length)
  // console.log(max(keysLength))
}

async function getExchanges() {
  const exchanges = await client.exchangeInfo()
  console.log(exchanges)
}

// client.ws.allTickers(tickers => {
//   console.log(tickers)
// })

// const str = 'BTCNEO'
// const match = str.match(/(.*)(BTC$)/)
// console.log(match)
