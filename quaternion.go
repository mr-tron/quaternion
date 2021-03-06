/*
Golang package implementing quaternion math
Purpose is to provide quaternion support under the MIT license as existing
Go quaternion packages are under more restrictive or unspecified licenses.

This project is licensed under the terms of the MIT license.
*/

package quaternion

import (
	"math"
)

// Quaternion represents a quaternion W+X*i+Y*j+Z*k
type Quaternion struct {
	W float64 // Scalar component
	X float64 // i component
	Y float64 // j component
	Z float64 // k component
}

// Conj returns the conjugate of a Quaternion (W,X,Y,Z) -> (W,-X,-Y,-Z)
func Conj(qin Quaternion) Quaternion {
	qout := Quaternion{}
	qout.W = +qin.W
	qout.X = -qin.X
	qout.Y = -qin.Y
	qout.Z = -qin.Z
	return qout
}

// Norm2 returns the L2-Norm of a Quaternion (W,X,Y,Z) -> W*W+X*X+Y*Y+Z*Z
func Norm2(qin Quaternion) float64 {
	return qin.W*qin.W + qin.X*qin.X + qin.Y*qin.Y + qin.Z*qin.Z
}

// Norm returns the L1-Norm of a Quaternion (W,X,Y,Z) -> Sqrt(W*W+X*X+Y*Y+Z*Z)
func Norm(qin Quaternion) float64 {
	return math.Sqrt(qin.W*qin.W + qin.X*qin.X + qin.Y*qin.Y + qin.Z*qin.Z)
}

// Scalar returns a scalar-only Quaternion representation of a float (W,0,0,0)
func Scalar(w float64) Quaternion {
	return Quaternion{W: w}
}

// Sum returns the vector sum of any number of Quaternions
func Sum(qin ...Quaternion) Quaternion {
	qout := Quaternion{}
	for _, q := range qin {
		qout.W += q.W
		qout.X += q.X
		qout.Y += q.Y
		qout.Z += q.Z
	}
	return qout
}

// Prod returns the non-commutative product of any number of Quaternions
func Prod(qin ...Quaternion) Quaternion {
	qout := Quaternion{1, 0, 0, 0}
	var w, x, y, z float64
	for _, q := range qin {
		w = qout.W*q.W - qout.X*q.X - qout.Y*q.Y - qout.Z*q.Z
		x = qout.W*q.X + qout.X*q.W + qout.Y*q.Z - qout.Z*q.Y
		y = qout.W*q.Y + qout.Y*q.W + qout.Z*q.X - qout.X*q.Z
		z = qout.W*q.Z + qout.Z*q.W + qout.X*q.Y - qout.Y*q.X
		qout = Quaternion{w, x, y, z}
	}
	return qout
}

// Unit returns the Quaternion rescaled to unit-L1-norm
func Unit(qin Quaternion) Quaternion {
	k := Norm(qin)
	return Quaternion{qin.W / k, qin.X / k, qin.Y / k, qin.Z / k}
}

// Inv returns the Quaternion conjugate rescaled so that Q Q* = 1
func Inv(qin Quaternion) Quaternion {
	k2 := Norm2(qin)
	q := Conj(qin)
	return Quaternion{q.W / k2, q.X / k2, q.Y / k2, q.Z / k2}
}

// Euler returns the Euler angles phi, theta, psi corresponding to a Quaternion
func Euler(q Quaternion) (float64, float64, float64) {
	r := Unit(q)
	phi := math.Atan2(2*(r.W*r.X+r.Y*r.Z), 1-2*(r.X*r.X+r.Y*r.Y))
	theta := math.Asin(2 * (r.W*r.Y - r.Z*r.X))
	psi := math.Atan2(2*(r.X*r.Y+r.W*r.Z), 1-2*(r.Y*r.Y+r.Z*r.Z))
	return phi, theta, psi
}

// FromEuler returns a Quaternion corresponding to Euler angles phi, theta, psi
func FromEuler(phi, theta, psi float64) Quaternion {
	q := Quaternion{}
	q.W = math.Cos(phi/2)*math.Cos(theta/2)*math.Cos(psi/2) +
		math.Sin(phi/2)*math.Sin(theta/2)*math.Sin(psi/2)
	q.X = math.Sin(phi/2)*math.Cos(theta/2)*math.Cos(psi/2) -
		math.Cos(phi/2)*math.Sin(theta/2)*math.Sin(psi/2)
	q.Y = math.Cos(phi/2)*math.Sin(theta/2)*math.Cos(psi/2) +
		math.Sin(phi/2)*math.Cos(theta/2)*math.Sin(psi/2)
	q.Z = math.Cos(phi/2)*math.Cos(theta/2)*math.Sin(psi/2) -
		math.Sin(phi/2)*math.Sin(theta/2)*math.Cos(psi/2)
	return q
}

// RotMat returns the rotation matrix (as float array) corresponding to a Quaternion
func RotMat(qin Quaternion) [3][3]float64 {
	q := Unit(qin)
	m := [3][3]float64{}
	m[0][0] = 1 - 2*(q.Y*q.Y+q.Z*q.Z)
	m[0][1] = 2 * (q.X*q.Y - q.W*q.Z)
	m[0][2] = 2 * (q.W*q.Y + q.X*q.Z)

	m[1][1] = 1 - 2*(q.Z*q.Z+q.X*q.X)
	m[1][2] = 2 * (q.Y*q.Z - q.W*q.X)
	m[1][0] = 2 * (q.W*q.Z + q.Y*q.X)

	m[2][2] = 1 - 2*(q.X*q.X+q.Y*q.Y)
	m[2][0] = 2 * (q.Z*q.X - q.W*q.Y)
	m[2][1] = 2 * (q.W*q.X + q.Z*q.Y)
	return m
}
