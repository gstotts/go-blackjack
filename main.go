package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Suit string
	Rank string
}

func (c Card) Int() []int {
	if c.Rank == "K" || c.Rank == "Q" || c.Rank == "J" {
		return []int{10}
	}

	if c.Rank == "A" {
		return []int{11, 1}
	}

	val, _ := strconv.Atoi(c.Rank)
	return []int{val}
}

type Deck struct {
	Cards    []Card
	Shuffled bool
}

type Hand []Card

func (h Hand) containsAces() (bool, int) {
	count := 0
	for _, card := range h {
		if card.Rank == "A" {
			count++
		}
	}

	if count == 0 {
		return false, 0
	}

	return true, count
}

func (h Hand) Total() int {
	var total int
	for _, card := range h {
		total += card.Int()[0]
	}

	if total <= 21 {
		return total
	}

	aces, count := h.containsAces()
	if total > 21 && aces {
		for i := 1; i <= count; i++ {
			total -= 10
			if total <= 21 {
				return total
			}
		}
	}
	return total
}

func generateID() string {
	return uuid.NewString()
}

func New_Deck() Deck {
	d := new(Deck)
	for _, suit := range []string{"Spades", "Hearts", "Clubs", "Diamonds"} {
		for _, rank := range []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"} {
			d.Cards = append(d.Cards, Card{
				Suit: suit,
				Rank: rank,
			})
		}
	}
	return *d
}

func (d *Deck) Shuffle() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })
	d.Shuffled = true
}

type Player struct {
	ID    string
	Name  string
	Hand  Hand
	Bet   float32
	Purse float32
}

type Game struct {
	ID           string
	Deck         Deck
	Player       Player
	Dealer_Shown Hand
	Dealer_Down  Card
}

func (g *Game) New() {
	g = new(Game)
	g.ID = generateID()
	g.Deck.Shuffled = false
	g.Deck = New_Deck()
	g.Deck.Shuffle()

	fmt.Println("Welcome to the Blackjack Table.  Good Luck!")
	fmt.Println("Your Name: ")
	var name string
	fmt.Scanln(&name)

	g.Player = Player{
		ID:    generateID(),
		Name:  name,
		Hand:  []Card{},
		Bet:   0,
		Purse: 1000,
	}

	fmt.Printf("\nGreat to meet you, %s!  Let's get playing...\n", g.Player.Name)
	g.Deal()
	g.Take_Bet()
	g.Show_Cards()
	g.Players_Action()
}

func (g *Game) Deal() {
	for i := 1; i <= 2; i++ {
		g.Player.Hand = append(g.Player.Hand, g.Deck.Cards[0])
		g.Dealer_Shown = []Card{g.Deck.Cards[1]}
		g.Dealer_Down = g.Deck.Cards[2]
		g.Deck.Cards = g.Deck.Cards[3:]
	}
}

func (g Game) Show_Cards() {
	fmt.Printf("\n\n   Your Were Dealt: %v, Total: %d\n", g.Player.Hand, g.Player.Hand.Total())
	fmt.Printf(" Dealer is Showing: %v\n", g.Dealer_Shown)
	fmt.Println("------------------------------------------------------")
	fmt.Printf("Your Current Wager: %.2f\n", g.Player.Bet)
}

func (g *Game) Take_Bet() {
	fmt.Printf("Your Current Cash: $%.2f\n\n", g.Player.Purse)
	fmt.Println("What would you like to wager? ")
	fmt.Scanln(&g.Player.Bet)
	for g.Player.Bet > g.Player.Purse || g.Player.Bet <= 0 {
		if g.Player.Bet > g.Player.Purse {
			fmt.Printf("Sorry.  That bet is too large.  You only have $%.2f available to bet.\n", g.Player.Purse)
		} else {
			fmt.Println("Sorry.  You cannot bet a value of $0.")
		}
		fmt.Println("What would you like to wager? ")
		fmt.Scanln(&g.Player.Bet)
	}
	g.Player.Purse = g.Player.Purse - g.Player.Bet
}

func (g *Game) Hit() {
	g.Player.Hand = append(g.Player.Hand, g.Deck.Cards[0])
	g.Deck.Cards = g.Deck.Cards[1:]
	for g.Player.Hand.Total() < 21 {
		g.Show_Cards()
		g.Players_Action()
	}

	g.Reveal_and_Settle()
}

func (g *Game) Reveal_and_Settle() {
	g.Dealer_Shown = append(g.Dealer_Shown, g.Dealer_Down)
	g.Show_Cards()
	if g.Player.Hand.Total() == 21 {
		fmt.Printf("Wahoo!  Way to go!  %s just won $%.2f!\n", g.Player.Name, (1.5 * g.Player.Bet))
		g.Player.Purse = g.Player.Purse + (2.5 * g.Player.Bet)
	} else if g.Player.Hand.Total() > 21 {
		fmt.Printf("Bust!  Better luck next time partner! %s just lost $%.2f\n", g.Player.Name, g.Player.Bet)
		g.Player.Purse = g.Player.Purse - g.Player.Bet
	} else {
		for g.Dealer_Shown.Total() < 17 {
			fmt.Println("Dealer chooses to hit.")
			g.Dealer_Shown = append(g.Dealer_Shown, g.Deck.Cards[0])
			g.Deck.Cards = g.Deck.Cards[1:]
			g.Show_Cards()
		}
		if g.Dealer_Shown.Total() > 21 {
			fmt.Printf("Dealer Busts!  %s just won $%.2f!\n", g.Player.Name, g.Player.Bet)
			g.Player.Purse = g.Player.Purse + (2 * g.Player.Bet)
		} else if g.Dealer_Shown.Total() == g.Player.Hand.Total() {
			fmt.Printf("Push!")
			g.Player.Purse = g.Player.Purse + g.Player.Bet
		} else if g.Dealer_Shown.Total() > g.Player.Hand.Total() {
			fmt.Printf("Dealer wins!  Better luck next time partner! %s just lost $%.2f\n", g.Player.Name, g.Player.Bet)
			g.Player.Purse = g.Player.Purse - g.Player.Bet
		} else {
			fmt.Printf("You beat the Dealer!  %s just won $%.2f\n", g.Player.Name, g.Player.Bet)
			g.Player.Purse = g.Player.Purse + (2 * g.Player.Bet)
		}
	}

	g.Player.Bet = 0
	fmt.Printf("Your Current Cash: $%.2f\n\n", g.Player.Purse)

	if g.Player.Purse <= 0 {
		fmt.Printf("Looks like you're out of money.  Try again some other time.  Thanks for playing %s!\n", g.Player.Name)
		os.Exit(0)
	}

	fmt.Println("Would you like to keep playing? y/n")
	var answer string
	fmt.Scanln(&answer)
	for !(strings.EqualFold(answer, "y")) && !(strings.EqualFold(answer, "n")) && !(strings.EqualFold(answer, "yes") && !(strings.EqualFold(answer, "no"))) {
		fmt.Println("Sorry, that is not a valid response.  Would you still like to play?")
		fmt.Scanln(&answer)
	}
	if strings.EqualFold(answer, "y") || strings.EqualFold(answer, "yes") {
		g.Next_Hand()
	} else {
		fmt.Printf("Until we meet again!  Good luck %s!\n", g.Player.Name)
		os.Exit(0)
	}
}

func (g *Game) Players_Action() {
	fmt.Println("Would you like to HIT or STAY?")
	var response string
	fmt.Scanln(&response)
	for !(strings.EqualFold("hit", response)) && !(strings.EqualFold("stay", response)) {
		fmt.Println("Sorry, I could not understand that.")
		fmt.Println("Would you like to HIT or STAY?")
		fmt.Scanln(&response)
	}
	if strings.EqualFold("hit", response) {
		g.Hit()
	} else if strings.EqualFold("stay", response) {
		g.Reveal_and_Settle()
	}
}

func (g *Game) Next_Hand() {
	g.Deck = New_Deck()
	g.Dealer_Down = Card{}
	g.Dealer_Shown = Hand{}
	g.Player.Hand = Hand{}
	g.Deck.Shuffle()
	g.Deal()
	g.Take_Bet()
	g.Show_Cards()
	g.Players_Action()
}

func main() {
	g := new(Game)
	g.New()
}
