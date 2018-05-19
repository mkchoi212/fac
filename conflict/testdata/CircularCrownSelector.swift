private func generateInitials() -> [String] {
  let randomString = UUID().uuidString
  let str = randomString.replacingOccurrences(of: "-", with: "")

  let abbrev = stride(from: 0, to: 18, by: 2).map { i -> String in
    let start = str.index(str.startIndex, offsetBy: i)
    let end = str.index(str.startIndex, offsetBy: i + 2)
    return String(str[start..<end])
  }

  return abbrev
}

<<<<<<< Updated upstream

private func randomColor() -> UIColor {
  let hue = ( CGFloat(arc4random() % 256) / 256.0 )               //  0.0 to 1.0
  let saturation = ( CGFloat(arc4random() % 128) / 256.0 ) + 0.5  //  0.5 to 1.0, away from white
  let brightness = ( CGFloat(arc4random() % 128) / 256.0 ) + 0.7  //  0.7 to 1.0, away from black
  return UIColor(hue: hue, saturation: saturation, brightness: brightness, alpha: 1)
}
||||||| merged common ancestors

private func randomColor() -> UIColor{
  let red:CGFloat = CGFloat(drand48())   
  let green:CGFloat = CGFloat(drand48()) 
  let blue:CGFloat = CGFloat(drand48())  
  return UIColor(red:red, green: green, blue: blue, alpha: 1.0)
} 
=======

private func randomColor() -> UIColor{
  let red = CGFloat(arc4random())   
  let green = CGFloat(arc4random()) 
  let blue = CGFloat(arc4random())  
  return UIColor(red:red, green: green, blue: blue, alpha: 1.0)
}
>>>>>>> Stashed changes
