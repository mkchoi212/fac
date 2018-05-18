//
//  CrownSelectorInterfaceController.swift
//  Circular Crown Selector WatchKit Extension
//
//  Created by Mike Choi on 11/24/17.
//  Copyright © 2017 Mike Choi. All rights reserved.
//

import WatchKit

class CrownSelectorInterfaceController: WKInterfaceController, WKCrownDelegate {
  @IBOutlet var circle1: WKInterfaceGroup!
  @IBOutlet var circle2: WKInterfaceGroup!
  @IBOutlet var circle3: WKInterfaceGroup!
  @IBOutlet var circle4: WKInterfaceGroup!
  @IBOutlet var circle5: WKInterfaceGroup!
  @IBOutlet var circle6: WKInterfaceGroup!
  @IBOutlet var circle7: WKInterfaceGroup!
  @IBOutlet var circle8: WKInterfaceGroup!
  @IBOutlet var circle9: WKInterfaceGroup!
  @IBOutlet var circle10: WKInterfaceGroup!
  @IBOutlet var circle11: WKInterfaceGroup!
  @IBOutlet var circle12: WKInterfaceGroup!
  @IBOutlet var currentLabel: WKInterfaceLabel!
  var circles : [WKInterfaceGroup]!

  var idx: Int!

  var deltaBuildUp: Double!
  let sensitivity = 0.2
  var abbrev : [String]!
  var fontColors : [UIColor] = []
 
  /// #6E4000
  let activeFontColor = #colorLiteral(red: 0.431372549, green: 0.2509803922, blue: 0, alpha: 1)
  /// #FF9403
  let activeColor = #colorLiteral(red: 1, green: 0.5803921569, blue: 0.01176470588, alpha: 1)
  /// #262628
  let inactiveColor = #colorLiteral(red: 0.1490196078, green: 0.1490196078, blue: 0.1568627451, alpha: 1)
  
  override func awake(withContext context: Any?) {
    super.awake(withContext: context)
    
    if let abbrev = context as? [String] {
      self.abbrev = abbrev
    } else {
      self.abbrev = generateInitials()
    }
    abbrev = fill(abbrev, with: "●")
    
    circles = [c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11, c12]

    _ = zip(circles, abbrev).map { (tup) -> Void in
      let (cir, str) = tup
      let fontColor = randomColor()
      fontColors.append(fontColor)
      cir.setBackgroundColor(inactiveColor)
      return cir.setBackgroundImage(stringToImage(str, color: fontColor))
    }
    c1.setBackgroundColor(activeColor)
  }
  
  override func willActivate() {
    // Make `crownSequncer` responsive
    crownSequencer.delegate = self
    crownSequencer.focus()
    deltaBuildUp = 0
    idx = 0
    setActive(0)
  }

  func fill(_ arr: [String], with str: String) -> [String] {
    let diff = 12 - arr.count
    let filledArray = arr + Array(repeating: str, count: diff)
    return filledArray
  }
  
  // Set group at idx active by changing color attributes
  func setActive(_ idx: Int) {
    circles[idx].setBackgroundColor(activeColor)
    circles[idx].setBackgroundImage(stringToImage(abbrev[idx], color: activeFontColor))
    currentLabel.setText(abbrev[idx])
  }
  
  func setInActive(_ idx: Int) {
    circles[idx].setBackgroundColor(inactiveColor)
    circles[idx].setBackgroundImage(stringToImage(abbrev[idx], color: fontColors[idx]))
  }
  
  // MARK: WKCrownDelegate
  func crownDidRotate(_ crownSequencer: WKCrownSequencer?, rotationalDelta: Double) {
    // Only act on crown rotation if `deltaBuildUp` is greater than sensitivity
    // for smoother / controllable scrolling experience
    deltaBuildUp = deltaBuildUp.sign != rotationalDelta.sign ? 0 : deltaBuildUp
    deltaBuildUp = deltaBuildUp + rotationalDelta
    
    if abs(deltaBuildUp) < sensitivity {
      return
    }

    setInActive(idx)

    idx = rotationalDelta > 0 ? idx + 1  : idx - 1;
    idx = idx % 12
    if idx < 0 {
      idx = 12 + idx
    }
    
    setActive(idx)
    deltaBuildUp = 0.0
  }

  // MARK: Helper Functions
  private func stringToImage(_ str: String, color: UIColor) -> UIImage? {
    let imageSize = CGSize(width: 23, height: 23)
    UIGraphicsBeginImageContextWithOptions(imageSize, false, 0)
    UIColor.clear.set()
    let rect = CGRect(origin: CGPoint.zero, size: imageSize)
    UIRectFill(rect)
   
    let style = NSMutableParagraphStyle()
    style.alignment = .center
    (str as NSString).draw(in: rect, withAttributes: [NSAttributedStringKey.font: UIFont.systemFont(ofSize: 13),
                                                      NSAttributedStringKey.paragraphStyle: style,
                                                      NSAttributedStringKey.foregroundColor: color,
                                                      NSAttributedStringKey.baselineOffset: -3.0])
    
    let image = UIGraphicsGetImageFromCurrentImageContext()
    UIGraphicsEndImageContext()
    return image
  }
 
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
}
